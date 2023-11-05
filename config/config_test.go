package config

/*
func TestTaggingAndGetRegion(t *testing.T) {
	tag := "testTag"

	st := Tagging(tag)

	if st == nil {
		t.Errorf("Tagging() failed, expected non-nil value, got nil")
	}

	stFromTable := GetRegion(tag)
	if stFromTable == nil {
		t.Errorf("GetRegion() failed, expected non-nil value, got nil")
	}

	if st != stFromTable {
		t.Errorf("GetRegion() did not return the expected singleTag instance")
	}
}

func TestSetAndGetBool(t *testing.T) {
	tag := "boolTag"
	st := Tagging(tag)

	expectedBool := true
	st.Set("boolKey", expectedBool)

	retrievedBool, err := st.GetBool("boolKey")
	if err != nil {
		t.Errorf("GetBool() failed: %v", err)
	}

	if retrievedBool != expectedBool {
		t.Errorf("GetBool() returned incorrect value, expected: %v, got: %v", expectedBool, retrievedBool)
	}
}

func TestSetAndGetString(t *testing.T) {
	tag := "stringTag"
	st := Tagging(tag)

	expectedString := "testString"
	st.Set("stringKey", expectedString)

	retrievedString, err := st.GetString("stringKey")
	if err != nil {
		t.Errorf("GetString() failed: %v", err)
	}

	if retrievedString != expectedString {
		t.Errorf("GetString() returned incorrect value, expected: %s, got: %s", expectedString, retrievedString)
	}
}

func TestSetAndGetNonExistentKey(t *testing.T) {
	tag := "nonExistentKeyTag"
	st := Tagging(tag)

	_, err := st.GetBool("nonExistentKey")
	if err == nil {
		t.Errorf("GetBool() on non-existent key did not return an error")
	}

	_, err = st.GetString("nonExistentKey")
	if err == nil {
		t.Errorf("GetString() on non-existent key did not return an error")
	}

	_, err = st.GetInt("nonExistentKey")
	if err == nil {
		t.Errorf("GetInt() on non-existent key did not return an error")
	}

	_, err = st.Get("nonExistentKey")
	if err == nil {
		t.Errorf("Get() on non-existent key did not return an error")
	}
}

func BenchmarkValue(b *testing.B) {
	g := Global()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Set(strconv.Itoa(i), rand.Intn(1000))
	}
	for i := 0; i < b.N; i++ {
		g.GetInt(strconv.Itoa(i))
	}
}
*/
