// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试创建新的favicon
func TestNewFavicon(t *testing.T) {
	favicon := NewFavicon()
	assert.NotNil(t, favicon, "favicon不应为空")
	assert.NotNil(t, favicon.body, "favicon body不应为空")
}

// 测试favicon的ServeHTTP方法
func TestFavicon_ServeHTTP(t *testing.T) {
	favicon := NewFavicon()

	// 创建测试请求
	req := httptest.NewRequest("GET", "/favicon.ico", nil)
	w := httptest.NewRecorder()

	// 调用ServeHTTP方法
	favicon.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code, "状态码应为200")
	assert.Equal(t, "image/x-icon", w.Header().Get("Content-Type"), "Content-Type应为image/x-icon")
	assert.NotEmpty(t, w.Body.Bytes(), "响应体不应为空")
}

// 测试favicon数据完整性
func TestFavicon_DataIntegrity(t *testing.T) {
	favicon := NewFavicon()

	// 验证favicon数据不为空且长度合理
	assert.NotEmpty(t, faviconBodys, "faviconBodys不应为空")
	assert.Greater(t, len(faviconBodys), 100, "favicon数据长度应大于100字节")

	// 验证两个实例的数据一致性
	favicon2 := NewFavicon()
	assert.Equal(t, favicon.body, favicon2.body, "两个favicon实例的数据应一致")
}

// 测试faviconBodys变量存在
func TestFaviconBodys_Exists(t *testing.T) {
	// 确保faviconBodys变量已定义且不为空
	assert.NotNil(t, faviconBodys, "faviconBodys变量应已定义")
	assert.NotEmpty(t, faviconBodys, "faviconBodys不应为空")
}
