package spamcontrol

import "testing"

func Test_getMajorityNumber(t *testing.T) {

	// act
	numbersAndExpectedResults := []struct {
		number         int
		expectedResult int
	}{{1, 1},
		{2, 2},
		{3, 2},
		{4, 3},
		{5, 3},
		{6, 4},
		{7, 4},
		{8, 5},
		{9, 5},
		{10, 6},
		{11, 6},
		{12, 7},
		{13, 7},
	}

	// arrange
	for _, entry := range numbersAndExpectedResults {

		// act
		result := getMajorityNumber(entry.number)

		// assert
		if result != entry.expectedResult {
			t.Fail()
			t.Logf("getMajorityNumber(%d) should have returned %q but returned %d instead.", entry.number, entry.expectedResult, result)
		}
	}

}

func Test_getMajorityStatus_AllStatusesAreError_ErrorStatusIsReturned(t *testing.T) {
	// arrange
	statuses := []Status{
		StatusError("1"),
		StatusError("2"),
		StatusError("3"),
	}

	// act
	exists, result := getMajorityStatus(statuses)

	// assert
	if !exists || result.Name() != StatusError("...").Name() {
		t.Fail()
		t.Logf("getMajorityStatus should have returned %q but returned %q instead.", StatusError("any"), result)
	}
}

func Test_getMajorityStatus_MajorityAvailable_MajorIsReturned(t *testing.T) {
	// arrange
	statuses := []Status{
		StatusNotInstalled(),
		StatusNotInstalled(),
		StatusNotInstalled(),
		StatusError("3"),
		StatusUnknown(),
	}

	// act
	exists, result := getMajorityStatus(statuses)

	// assert
	if !exists || result.Name() != StatusNotInstalled().Name() {
		t.Fail()
		t.Logf("getMajorityStatus should have returned %q but returned %q instead.", StatusNotInstalled(), result)
	}
}

func Test_getMajorityStatus_NoMajority_ResultIsFalse(t *testing.T) {
	// arrange
	statuses := []Status{
		StatusError("1"),
		StatusUnknown(),
		StatusUpToDate(),
	}

	// act
	exists, result := getMajorityStatus(statuses)

	// assert
	if exists || result != nil {
		t.Fail()
		t.Logf("getMajorityStatus returned a status even though there is no majority.")
	}
}

func Test_getMajorityStatus_NoStatuses_ResultIsFalse(t *testing.T) {
	// arrange
	statuses := []Status{}

	// act
	exists, result := getMajorityStatus(statuses)

	// assert
	if exists || result != nil {
		t.Fail()
		t.Logf("getMajorityStatus returned a status even though no statuses were given.")
	}
}

func Test_getMajorityStatus_Nil_ResultIsFalse(t *testing.T) {
	// act
	exists, result := getMajorityStatus(nil)

	// assert
	if exists || result != nil {
		t.Fail()
		t.Logf("getMajorityStatus returned a status even though no statuses were given.")
	}
}
