package slice

import "testing"

func TestContainsLowerExclude(t *testing.T) {
	vIncludes, vNewArgs := ContainsLowerExclude([]string{"aaa", "bbb", "ccc"}, nil)
	if vIncludes {
		t.Error("Excepcted includes to be false, received ", vIncludes)
	}
	if len(vNewArgs) != 3 {
		t.Error("Expected 3 newArgs, received ", len(vNewArgs))
	}
	vIncludes, vNewArgs = ContainsLowerExclude([]string{"aaa", "bbb", "ccc"}, []string{"aaa"})
	if !vIncludes {
		t.Error("Excepcted includes to be true, received ", vIncludes)
	}
	if len(vNewArgs) != 2 {
		t.Error("Expected 2 newArgs, received ", len(vNewArgs))
	}
	vIncludes, vNewArgs = ContainsLowerExclude([]string{"aaa", "bbb", "ccc"}, []string{"eee"})
	if vIncludes {
		t.Error("Excepcted includes to be false, received ", vIncludes)
	}
	if len(vNewArgs) != 3 {
		t.Error("Expected 3 newArgs, received ", len(vNewArgs))
	}

}
