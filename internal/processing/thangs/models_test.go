package thangs

import "testing"

func TestModelIdFromUrlOnDesignerUrl(t *testing.T) {
	url := "/designer/foo"
	expected := ""
	actual := modelIdFromUrl(url)
	if actual != expected {
		t.Errorf("expected [%s] but got [%+v]", expected, actual)
	}
}

func TestModelIdFromUrlOnModelUrl(t *testing.T) {
	url := "/designer/foo/3d-model/bar-1234567"
	expected := "1234567"
	actual := modelIdFromUrl(url)
	if actual != expected {
		t.Errorf("expected [%s] but got [%+v]", expected, actual)
	}
}
