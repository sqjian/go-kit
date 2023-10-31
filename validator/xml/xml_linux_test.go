package xml_test

//func TestValidateXmlWithXsd(t *testing.T) {
//	openFile := func(fileName string) []byte {
//		fileContent, fileContentErr := ioutil.ReadFile(fileName)
//		if fileContentErr != nil {
//			t.Fatal(fileContentErr)
//		}
//		return fileContent
//	}
//
//	type args struct {
//		xmlData []byte
//		xsdData []byte
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		{
//			"test1",
//			args{
//				openFile("test.xml"),
//				openFile("test.xsd"),
//			},
//			false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if err := ValidateXmlWithXsd(tt.args.xmlData, tt.args.xsdData); (err != nil) != tt.wantErr {
//				t.Errorf("ValidateXmlWithXsd() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
