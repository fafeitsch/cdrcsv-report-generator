package samples

import (
	"bufio"
	"bytes"
	"testing"
)

func TestCreate(t *testing.T) {
	options := Options{Count: 10}
	writer := new(bytes.Buffer)
	err := Create(&options, writer)
	if err != nil {
		t.Errorf("Creating threw unexpected error: %v", err)
	}
	reader := bufio.NewReader(writer)
	expectedLine := "\"\",\"4711\",\"0815\",\"from_public\",\"\"John Doe\" <4711>\",\"SIP/from_public-0000012\",\"SIP/deskphone_of_boss_000015a\",\"DIAL\",\"SIP/deskphone_of_boss\",\"2019-05-09 08:06:11\",\"2019-05-09 08:06:30\",\"2019-05-09 08:30:12\",1441,1422,\"ANSWERED\",\"DOCUMENTATION\",\"\",\"2265436.50\"\n"
	for i := 0; i < options.Count; i++ {
		line, _ := reader.ReadString('\n')
		if line == "" {
			t.Errorf("Line %d is empty", i)
		} else if line != expectedLine {
			t.Errorf("Line %d differs (expected vs. actual):\n%s\n%s", i, expectedLine, line)
		}
	}
}
