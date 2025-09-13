package thangs

import "testing"

func TestUrlToMetaOnEmptyString(t *testing.T) {
	url := ""
	expected := urlMeta{}
	actual, err := urlToMeta(url)
	if err == nil {
		t.Errorf("expected an error but got [%+v]", err)
	}
	if actual != expected {
		t.Errorf("expected [%s] but got [%+v]", expected, actual)
	}
}

func TestUrlToMetaOnDesignerUrl(t *testing.T) {
	url := "/designer/foo"
	expected := urlMeta{}
	actual, err := urlToMeta(url)
	if err == nil {
		t.Errorf("expected an error but got [%+v]", err)
	}
	if actual != expected {
		t.Errorf("expected [%s] but got [%+v]", expected, actual)
	}
}

func TestUrlToMetaOnModelUrl(t *testing.T) {
	url := "/designer/foo/3d-model/Foo%20Bar-1234567"
	expected := urlMeta{id:"1234567", title: "Foo Bar"}
	actual, err := urlToMeta(url)
	if err != nil {
		t.Error(err)
	}
	if actual != expected {
		t.Errorf("expected [%s] but got [%+v]", expected, actual)
	}
}

func TestUrlToMetaOnModelUrlWithQuery(t *testing.T) {
	url := "/designer/foo/3d-model/bar-1234567?image=12345"
	expected := urlMeta{id:"1234567", title: "bar"}
	actual, err := urlToMeta(url)
	if err != nil {
		t.Error(err)
	}
	if actual != expected {
		t.Errorf("expected [%s] but got [%+v]", expected, actual)
	}
}
