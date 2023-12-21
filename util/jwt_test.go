package util_test

import (
	"swpr/util"
	"testing"
)

func TestJwtCreateToken(t *testing.T) {
	type args struct {
		userID int64
		jwtKey string
		ttl    string
	}
	tests := []struct {
		name            string
		args            args
		wantTokenString string
		wantErr         bool
	}{
		{
			name: "Success create jwt",
			args: args{
				userID: 1,
				jwtKey: "Test123TestKeyJwt",
				ttl:    "1m",
			},
			wantTokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk3MDMxNTcxMDcsImlhdCI6MTcwMzE1NzA0Nywic3ViIjoxfQ.OKVPwew56tl1nnUlfwVcus60FvGy4NnHf3Bk_3DpAaw",
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTokenString, err := util.JwtCreateToken(tt.args.userID, tt.args.jwtKey, tt.args.ttl)
			if (err != nil) != tt.wantErr {
				t.Errorf("JwtCreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			//checking Claim fields
			wantJwtToken, _ := util.JwtParseToken(tt.wantTokenString, tt.args.jwtKey)
			actualJwtToken, _ := util.JwtParseToken(gotTokenString, tt.args.jwtKey)

			wantSub, _ := wantJwtToken.Claims.GetSubject()
			actualSub, _ := wantJwtToken.Claims.GetSubject()

			if wantSub != actualSub {
				t.Errorf("JwtCreateToken() claim actual = %v, want %v", actualJwtToken.Claims, wantJwtToken.Claims)
			}

			//checking token valid
			isValid, _ := util.JwtVerify(tt.wantTokenString, tt.args.jwtKey)
			if !isValid {
				t.Error("JwtCreateToken() Invalid")
			}
		})
	}
}

func TestJwtVerify(t *testing.T) {
	type args struct {
		tokenString string
		jwtKey      string
	}
	tests := []struct {
		name        string
		args        args
		wantIsValid bool
		wantErr     bool
	}{
		{
			name: "Success",
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk3MDMxNTcxMDcsImlhdCI6MTcwMzE1NzA0Nywic3ViIjoxfQ.OKVPwew56tl1nnUlfwVcus60FvGy4NnHf3Bk_3DpAaw",
				jwtKey:      "Test123TestKeyJwt",
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Error invalid key",
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk3MDMxNTcxMDcsImlhdCI6MTcwMzE1NzA0Nywic3ViIjoxfQ.OKVPwew56tl1nnUlfwVcus60FvGy4NnHf3Bk_3DpAaw",
				jwtKey:      "Test123TestKeyJwts",
			},
			wantIsValid: false,
			wantErr:     false,
		},
		{
			name: "Error expired",
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDMxNTcxMDcsImlhdCI6MTcwMzE1NzA0Nywic3ViIjoxfQ.uBc-Ok5i3I5W-nLWyFUSm5ODOmrgvXpzBtPkworY1no",
				jwtKey:      "Test123TestKeyJwt",
			},
			wantIsValid: false,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIsValid, _ := util.JwtVerify(tt.args.tokenString, tt.args.jwtKey)
			if gotIsValid != tt.wantIsValid {
				t.Errorf("JwtVerify() gotIsValid = %v, want %v", gotIsValid, tt.wantIsValid)
			}
		})
	}
}
