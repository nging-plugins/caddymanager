/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package s3client

import (
	"crypto/tls"
	"net/http"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/s3manager"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func Connect(m *dbschema.NgingCloudStorage) (client *minio.Client, err error) {
	isSecure := m.Secure == `Y`
	options := &minio.Options{
		Creds:  credentials.NewStaticV4(m.Key, m.Secret, ""),
		Secure: isSecure,
		Region: m.Region,
	}
	if isSecure {
		options.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	client, err = minio.New(m.Endpoint, options)
	return
}

func New(m *dbschema.NgingCloudStorage, editableMaxSize int) *s3manager.S3Manager {
	return s3manager.New(Connect, m, editableMaxSize)
}
