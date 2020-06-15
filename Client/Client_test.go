package main

import "testing"

func TestClient(t *testing.T) {
	cli := NewClient()
	value, err := cli.Read("2")
	t.Logf("First read: %v, error: %v\n", value, err)
	if err := cli.Put("1", "value1"); err != nil{
		t.Fatal(err)
	}
	if err := cli.Put("2", "value2"); err != nil{
		t.Fatal(err)
	}
	if err := cli.Put("2", "value2.2"); err != nil{
		t.Fatal(err)
	}

	if value, err := cli.Read("2"); err != nil || value != "value2.2"{
		t.Fatal(err)
	}
	if err := cli.Delete("2"); err != nil{
		t.Fatal(err)
	}
	if value, err := cli.Read("2"); err == nil || value != ""{
		t.Fatal("Read deleted value")
	}

	t.Logf("Last read: %v, error: %v\n", value, err)
}