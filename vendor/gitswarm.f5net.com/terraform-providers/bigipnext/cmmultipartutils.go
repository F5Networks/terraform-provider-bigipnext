/*
Copyright 2024 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
// Package bigipnext interacts with BIGIP-NEXT/CM systems using the OPEN API.
package bigipnext

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
)

func WriteFile(w *multipart.Writer, fieldname string, file *os.File) error {
	p, err := PreparePart(w, fieldname, filepath.Base(file.Name()))
	if err != nil {
		return err
	}
	_, err = io.Copy(p, file)
	return err
}
func WriteField(w *multipart.Writer, fieldname, value string) error {
	p, err := PreparePart(w, fieldname, "")
	if err != nil {
		return err
	}
	_, err = p.Write([]byte(value))
	return err
}
func PreparePart(w *multipart.Writer, fieldname, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	// h.Set(echo.HeaderContentDisposition, fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldname, filename))
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldname, filename))
	// h.Set(echo.HeaderContentType, echo.MIMETextPlain)
	h.Set("Content-Type", "text/plain")
	return w.CreatePart(h)
}

// func VerifyHalContent(header http.Header) error {
// 	contentType := header.Get("Content-Type")
// 	if contentType != utils.HalJSONHeader {
// 		return fmt.Errorf("contentType is empty or not equal to HalJson")
// 	}
// 	return nil
// }

// func VerifyHalLinksPlural(hal waf.HalLinksPlural) error {
// 	if hal == (waf.HalLinksPlural{}) || hal.Links == nil || hal.Links.Self == nil || hal.Links.Self.Href == "" {
// 		return fmt.Errorf("Hal links are empty")
// 	}
// 	return nil
// }
