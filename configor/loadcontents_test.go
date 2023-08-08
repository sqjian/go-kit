package configor

import "testing"

func TestLoadJsonContents(t *testing.T) {
	type args struct {
		obj  any
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				obj: struct {
				}{},
				data: nil,
			},
			wantErr: false,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadJsonContents(tt.args.obj, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("LoadJsonContents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
