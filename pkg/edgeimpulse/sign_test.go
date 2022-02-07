package edgeimpulse

import (
	"fmt"
	"testing"
)

func Test_sign(t *testing.T) {
	type args struct {
		key  string
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "Test sign", args: args{key: "olia", data: []byte("data")}, wantErr: false,
			want: fmt.Sprintf(jsonPayload, "f4cea4f96ec79db63d82db2eb4e34f59d48b29e37ca6b291d8662e762809418f",
				4, "f375c4e34ab05c45e3744d92cdfb3dd5f1244e54afcd59d9fbe24a64e7b1214c")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sign(tt.args.key, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("sign() = %v, want %v", got, tt.want)
			}
		})
	}
}
