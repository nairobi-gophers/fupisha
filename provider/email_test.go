package provider

import (
	"bytes"
	"testing"
	"time"
)

func TestParseTemplates(t *testing.T) {
	tpl, err := parseTemplates("./testdata")
	if err != nil {
		t.Fatal(err)
	}

	testTemplateName := "hello"
	testTemplateContent := struct {
		SiteURL            string
		SiteName           string
		VerificationURL    string
		VerificationExpiry time.Time
	}{
		SiteURL:            "https://fupisha.io",
		SiteName:           "Fupisha",
		VerificationURL:    "https://fupisha.io/verify/?v=123456789",
		VerificationExpiry: time.Now().Add(time.Minute * 15),
	}

	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, testTemplateName, testTemplateContent); err != nil {
		t.Fatal(err)
	}

	want := []byte(`
	<html>
	<head>
	<title>Fupisha</title>
	</head>
	<body>
    <h1>Hello, </h1>
		<p>
		Thank you for signing up for Fupisha as an early user! Please verify your email 
		<a href="https://fupisha.io/verify/?v=123456789">
		<button>Verify</button>
		</a>
		</p>

		<p> The link above expires in 14 minutes</p>

		<a href="https://fupisha.io">The Fupisha Team</a>
	</body>
	</html>`)

	if bytes.Equal(buf.Bytes(), want) {
		t.Fatalf("got %q; want %q", buf.String(), string(want))
	}

}
