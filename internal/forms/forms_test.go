package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Vaild(t *testing.T) {
	r := httptest.NewRequest("POST", "/whaterver", nil)
	form := New(r.PostForm)
	if !form.Valid() {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whaterver", nil)
	form := New(r.PostForm)
	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}
	postedData := url.Values{}
	postedData.Add("a", "aa")
	postedData.Add("b", "bb")
	postedData.Add("c", "cc")

	r = httptest.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData

	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}

}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	if form.Has("a") {
		t.Error("form shows has field when it does not")
	}
	postedData := url.Values{}
	postedData.Add("a", "aa")
	form = New(postedData)

	if !form.Has("a") {
		t.Error("shows form does not have field when it should")
	}

}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	form.MinLength("a", 3)
	if form.Valid() {
		t.Error("form shows min length met for non-existent field")
	}

	isError := form.Errors.Get("a")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	postedData := url.Values{}
	postedData.Add("a", "1234")
	form = New(postedData)

	form.MinLength("a", 100)
	if form.Valid() {
		t.Error("form shows min length met when data is shorter")
	}

	postedData = url.Values{}
	postedData.Add("a", "1234abc")
	form = New(postedData)

	form.MinLength("a", 1)
	if !form.Valid() {
		t.Error("form shows min length is not met when it is")
	}

	isError = form.Errors.Get("a")
	if isError != "" {
		t.Error("should not have an error, but got one")
	}

}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)
	form.IsEmail("a")

	if form.Valid() {
		t.Error("form shows valid email for non-existent field")
	}

	postedData = url.Values{}
	postedData.Add("a", "1234")
	form = New(postedData)
	form.IsEmail("a")

	if form.Valid() {
		t.Error("got valid for invalid email address")
	}

	postedData = url.Values{}
	postedData.Add("b", "1234@qq.com")
	form = New(postedData)
	form.IsEmail("b")

	if !form.Valid() {
		t.Error("got an invalid email when we should not have")
	}
}
