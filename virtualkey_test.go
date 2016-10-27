// Go FIDO U2F Library
// Copyright 2015 The Go FIDO U2F Library Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package u2f

import (
    "testing"
    "fmt"
)

func TestVirtualKey(t *testing.T) {
    
    vk, err := NewVirtualKey()
    if err != nil {
        t.Error(err)
        t.FailNow()
    }

    app_id := "http://localhost"

    // Generate registration request
    c1, _ := NewChallenge(app_id, []string{app_id})
    //fmt.Printf("Challenge: %+v\n", c)
    registerReq := c1.RegisterRequest()

    // Pass to virtual token
    resp, err := vk.HandleRegisterRequest(*registerReq)
    if err != nil {
        t.Error(err)
        t.FailNow()
    }

    // Register virtual token
    // Attestation cert is self signed, so skip checking that
    reg, err := Register(*resp, *c1, &Config{SkipAttestationVerify: true})
    if err != nil {
        t.Error(err)
        t.FailNow()
    }

    // Send authentication request to the browser / token.
    c2, _ := NewChallenge(app_id, []string{app_id})
    signReq := c2.SignRequest(*reg)

    // Pass to virtual token
    signResp, err := vk.HandleAuthenticationRequest(*signReq)
    if err != nil {
        t.Error(err)
        t.FailNow()
    }

    // Read response from the browser / token.
    newCounter, err := reg.Authenticate(*signResp, *c2, 0)
    if err != nil {
        t.Error(err)
        t.FailNow()
    }

    fmt.Printf("Counter: %+v\n", newCounter)

}
