package main

import (
	"testing"
)

func TestMaxOf(t *testing.T) {
  t1 := []string{"1", "2", "3", "4", "5"}
  r1,err := MaxOf(t1)
  if err != nil { t.Fail() }
  if r1 != 5 { t.Fail() }

  t2 := []string{"66", "0", "12", "0"}
  r2,err := MaxOf(t2)
  if err != nil { t.Fail() }
  if r2 != 66 { t.Fail() }
}

func TestCompareVersion(t *testing.T){
  r1,err := CompareVersion("1.0.0", "1.0.0")
  if err != nil || r1 != 0 { t.Fail() }

  r2,err := CompareVersion("1.0.1", "1.0.0")
  if err != nil || r2 != 1 { t.Fail() }

  r3,err := CompareVersion("1.0.0", "1.0.1")
  if err != nil || r3 != -1 { t.Fail() }

  r4,err := CompareVersion("1.0", "1.0.0")
  if err != nil || r4 != 0 { t.Fail() }

  r5,err := CompareVersion("2.0", "1.99.55")
  if err != nil || r5 != 1 { t.Fail() }

  r6,err := CompareVersion("1.99.55", "2.0")
  if err != nil || r6 != -1 { t.Fail() }

  r7,err := CompareVersion("1", "0.456")
  if err != nil || r7 != 1 { t.Fail() }
}
