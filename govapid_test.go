package govapid

import (
	"encoding/base64"
	"fmt"
	"testing"
)

const succeed = "\u2705"
const failed = "\u274C"

//TestGenerateVAPID will test vapidkeys generation
func TestGenerateVAPID(t *testing.T) {
	t.Log("Make sure correct Public and Private VAPID keys are generated")
	{
		VAPIDKeys, err := GenerateVAPID()
		if err != nil {
			t.Fatalf("\t%s\tShould generate VAPID keys, got an error %s", failed, err.Error())
		}
		t.Logf("\t%s\tPublic key and Private key generated successfully", succeed)

		//Validate keys length
		if len(VAPIDKeys.Public) != 87 {
			t.Fatalf("\t%s\tInvalid Public VAPID key", failed)
		}
		t.Logf("\t%s\tPrivate key length is valid", succeed)

		if len(VAPIDKeys.Private) != 43 {
			t.Fatalf("\t%s\tInvalid Private VAPID key", failed)
		}
		t.Logf("\t%s\tPublic key length is valid", succeed)

		//Test base64 decoding
		_, err = base64.RawURLEncoding.DecodeString(VAPIDKeys.Private)
		if err != nil {
			t.Fatalf("\t%s\tInvalid Private key: %s", failed, err.Error())
		}
		t.Logf("\t%s\tPrivate key is valid", succeed)

		_, err = base64.RawURLEncoding.DecodeString(VAPIDKeys.Public)
		if err != nil {
			t.Fatalf("\t%s\tInvalid Public key: %s", failed, err.Error())
		}
		t.Logf("\t%s\tPublic key is valid", succeed)
	}
}

func ExampleGenerateVAPID() {
	VAPIDkeys, err := GenerateVAPID()
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(VAPIDkeys.Public, VAPIDkeys.Private)
	fmt.Println(len(VAPIDkeys.Public), len(VAPIDkeys.Private))
	// Output:
	// 87 43
}
