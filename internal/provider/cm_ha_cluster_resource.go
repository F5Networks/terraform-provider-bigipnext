package provider

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"os/exec"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	bigipnextsdk "gitswarm.f5net.com/terraform-providers/bigipnext"
)

var (
	_ resource.Resource = &NextCMHAClusterResource{}
	// _ resource.ResourceWithImportState = &NextCMHAClusterResource{}
)

func NewNextCMHAClusterResource() resource.Resource {
	return &NextCMHAClusterResource{}
}

type NextCMHAClusterResource struct {
	client *bigipnextsdk.BigipNextCM
}

type Nodes struct {
	NodeIP      types.String `tfsdk:"node_ip"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	Fingerprint types.String `tfsdk:"fingerprint"`
}

type NextCMHAClusterResourceModel struct {
	Nodes       []Nodes      `tfsdk:"nodes"`
	AgentNodes  types.List   `tfsdk:"agent_nodes"`
	ServerNodes types.List   `tfsdk:"server_nodes"`
	ID          types.String `tfsdk:"id"`
}

func (r *NextCMHAClusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_ha_cluster"
}

func (r *NextCMHAClusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a HA Cluster of BIG-IP Next Central Manager instances",
		Attributes: map[string]schema.Attribute{
			"nodes": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"node_ip": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "IP address of the node that will be added to the cluster",
						},
						"username": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The username of the node",
						},
						"password": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The password of the node",
							Sensitive:           true,
						},
						"fingerprint": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The fingerprint of the node in the SHA256 format",
						},
					},
				},
				Required: true,
			},
			"agent_nodes": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "List of nodes that are marked as agent nodes",
			},
			"server_nodes": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "List of nodes that are marked as control plane nodes",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *NextCMHAClusterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (r *NextCMHAClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resCfg *NextCMHAClusterResourceModel
	var nodes []bigipnextsdk.CMHANodes
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var nodeCheck []string
	for _, node := range resCfg.Nodes {
		if node.Fingerprint.IsNull() {
			fingerprint, err := getFingerPrint(node.NodeIP.String())
			if err != nil {
				resp.Diagnostics.AddError(fmt.Sprintf("error getting fingerprint of the node %s: ", node.NodeIP), err.Error())
				return
			}
			node.Fingerprint = types.StringValue(fingerprint)
		}

		n := bigipnextsdk.CMHANodes{
			NodeAddress: node.NodeIP.ValueString(),
			Username:    node.Username.ValueString(),
			Password:    node.Password.ValueString(),
			Fingerprint: node.Fingerprint.ValueString(),
		}

		nodeCheck = append(nodeCheck, n.NodeAddress)
		nodes = append(nodes, n)
	}

	res, err := r.client.CreateCMHACluster(nodes)

	if err != nil {
		resp.Diagnostics.AddError("error creating CM HA cluster", err.Error())
		return
	}

	log.Printf("Started CM HA Cluster creation: %v", res)

	res2, err := r.client.CheckCMHANodesStatus(nodeCheck)
	if err != nil {
		resp.Diagnostics.AddError("error creating CM HA cluster", err.Error())
		return
	}

	serverNodes, agentNodes := getServerAndAgentNodes(res2)
	cmServerIP := extractIPFromUrl(r.client.Host)
	id := fmt.Sprintf("central-manager-server-%s", cmServerIP)

	resp.State.SetAttribute(ctx, path.Root("server_nodes"), serverNodes)
	resp.State.SetAttribute(ctx, path.Root("agent_nodes"), agentNodes)
	resp.State.SetAttribute(ctx, path.Root("id"), id)
	resp.State.SetAttribute(ctx, path.Root("nodes"), resCfg.Nodes)
}

func (r *NextCMHAClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateCfg *NextCMHAClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetCMHANodes()
	if err != nil {
		resp.Diagnostics.AddError("error reading CM HA cluster", err.Error())
		return
	}

	serverNodes, agentNodes := getServerAndAgentNodes(res)
	cmServerIP := extractIPFromUrl(r.client.Host)
	id := fmt.Sprintf("central-manager-server-%s", cmServerIP)

	resp.State.SetAttribute(ctx, path.Root("server_nodes"), serverNodes)
	resp.State.SetAttribute(ctx, path.Root("agent_nodes"), agentNodes)
	resp.State.SetAttribute(ctx, path.Root("id"), id)
	resp.State.SetAttribute(ctx, path.Root("nodes"), stateCfg.Nodes)
}

func (r *NextCMHAClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var PlanCfg *NextCMHAClusterResourceModel
	var StateCfg *NextCMHAClusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &PlanCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &StateCfg)...)

	planNodes := getNodeIPs(PlanCfg.Nodes)
	stateNodes := getNodeIPs(StateCfg.Nodes)

	deleteNodes := nodesDiference(stateNodes, planNodes)
	addNodes := nodesDiference(planNodes, stateNodes)
	var nodeCheck []string
	if slices.Compare(planNodes, stateNodes) != 0 {

		if len(deleteNodes) > 0 {
			tflog.Info(ctx, fmt.Sprintf("deleting nodes: %v", deleteNodes))
			r.client.DeleteCMHANodes(deleteNodes)
		}

		if len(addNodes) > 0 {
			tflog.Info(ctx, fmt.Sprintf("adding nodes: %v", addNodes))

			var nodes []bigipnextsdk.CMHANodes
			for _, node := range PlanCfg.Nodes {
				doAddOperation := false
				for _, addNode := range addNodes {
					if node.NodeIP.String() == addNode {
						doAddOperation = true
						break
					}
				}
				if !doAddOperation {
					continue
				}
				if node.Fingerprint.IsNull() {
					fingerprint, err := getFingerPrint(node.NodeIP.String())
					if err != nil {
						resp.Diagnostics.AddError(fmt.Sprintf("error getting fingerprint of the node %s: ", node.NodeIP), err.Error())
						return
					}
					node.Fingerprint = types.StringValue(fingerprint)
				}

				n := bigipnextsdk.CMHANodes{
					NodeAddress: node.NodeIP.ValueString(),
					Username:    node.Username.ValueString(),
					Password:    node.Password.ValueString(),
					Fingerprint: node.Fingerprint.ValueString(),
				}

				nodeCheck = append(nodeCheck, n.NodeAddress)
				nodes = append(nodes, n)
			}

			_, err := r.client.CreateCMHACluster(nodes)
			if err != nil {
				resp.Diagnostics.AddError("error updating CM HA cluster", err.Error())
				return
			}
		}
	}

	res, err := r.client.CheckCMHANodesStatus(nodeCheck)
	if err != nil {
		resp.Diagnostics.AddError("error reading CM HA cluster", err.Error())
		return
	}

	serverNodes, agentNodes := getServerAndAgentNodes(res)
	cmServerIP := extractIPFromUrl(r.client.Host)
	id := fmt.Sprintf("central-manager-server-%s", cmServerIP)

	resp.State.SetAttribute(ctx, path.Root("server_nodes"), serverNodes)
	resp.State.SetAttribute(ctx, path.Root("agent_nodes"), agentNodes)
	resp.State.SetAttribute(ctx, path.Root("id"), id)
	resp.State.SetAttribute(ctx, path.Root("nodes"), PlanCfg.Nodes)
}

func (r *NextCMHAClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var stateCfg *NextCMHAClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nodes := getNodeIPs(stateCfg.Nodes)
	r.client.DeleteCMHANodes(nodes)

	res, _ := r.client.GetCMHANodes()

	nodeCount := 0
	for _, n := range res {
		if n.Status.Ready {
			nodeCount++
		}
	}

	if nodeCount > 1 {
		tflog.Error(ctx, "unable to delete all the nodes of the cluster")
		resp.Diagnostics.AddError("delete operation failed", "unable to delete all the nodes of the cluster")
	}

	stateCfg.ID = types.StringValue("")
}

func nodesDiference(nodes1, nodes2 []string) []string {
	var diff []string
	for _, node1 := range nodes1 {
		c := node1
		for _, node2 := range nodes2 {
			if node1 == node2 {
				c = ""
			}
		}
		if c != "" {
			diff = append(diff, c)
		}
	}
	return diff
}

func getNodeIPs(nodes []Nodes) []string {
	var nodeIPs []string
	for _, node := range nodes {
		nodeIPs = append(nodeIPs, node.NodeIP.String())
	}
	return nodeIPs
}

func getServerAndAgentNodes(res []bigipnextsdk.CMHANodesStatus) ([]string, []string) {
	var serverNodes []string
	var agentNodes []string

	for _, node := range res {
		if node.Spec.NodeType == "agent" {
			agentNodes = append(agentNodes, node.Spec.NodeAddress)
		} else if node.Spec.NodeType == "server" {
			serverNodes = append(serverNodes, node.Spec.NodeAddress)
		}
	}

	return serverNodes, agentNodes
}

func getFingerPrint(ip string) (string, error) {
	cmd := fmt.Sprintf(`openssl s_client -connect %s:443 < /dev/null 2>/dev/null | openssl x509 -fingerprint -noout -in /dev/stdin -SHA256`, ip)
	res := exec.Command("bash", "-c", cmd)
	out, err := res.Output()

	if err != nil {
		return "", fmt.Errorf("error getting fingerprint of the node %s: %v", ip, err)
	}

	o := strings.Split(string(out), "=")
	fingerprint := strings.ToLower(
		strings.TrimSpace(
			strings.ReplaceAll(o[1], ":", ""),
		),
	)
	return string(fingerprint), nil
}

func extractIPFromUrl(urlStr string) string {
	parsedURL, _ := url.Parse(urlStr)
	host, _, _ := net.SplitHostPort(parsedURL.Host)
	ip := net.ParseIP(host)
	return ip.String()
}
