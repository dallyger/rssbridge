package shopware

import (
	"testing"
)

func TestPluginChangelogTitle(t *testing.T) {
	expected := "1.2.3 | 6.1.0 - 6.3.1.1 | 21 September 2020"
	actual := pluginChangelogTitle(`
                        1.2.3




                                    6.1.0 - 6.3.1.1


                                21 September 2020
		`)
	if actual != expected {
		t.Errorf("expected: [%s], got [%s]", expected, actual)
	}
}
func TestPluginChangelogId(t *testing.T) {
	expected := "1.2.3"
	actual := pluginChangelogId(`
                        1.2.3




                                    6.1.0 - 6.3.1.1


                                21 September 2020
		`)
	if actual != expected {
		t.Errorf("expected: [%s], got [%s]", expected, actual)
	}
}
