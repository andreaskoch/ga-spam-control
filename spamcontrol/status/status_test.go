package status

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
		ErrorDefault,
		ErrorDefault,
		ErrorDefault,
	}

	// act
	exists, result := getMajorityStatus(statuses)

	// assert
	if !exists || result.Name() != ErrorDefault.Name() {
		t.Fail()
		t.Logf("getMajorityStatus should have returned %q but returned %q instead.", ErrorDefault, result)
	}
}

func Test_getMajorityStatus_MajorityAvailable_MajorIsReturned(t *testing.T) {
	// arrange
	statuses := []Status{
		NotInstalled,
		NotInstalled,
		NotInstalled,
		ErrorDefault,
		Unknown,
	}

	// act
	exists, result := getMajorityStatus(statuses)

	// assert
	if !exists || result.Name() != NotInstalled.Name() {
		t.Fail()
		t.Logf("getMajorityStatus should have returned %q but returned %q instead.", NotInstalled, result)
	}
}

func Test_getMajorityStatus_NoMajority_ResultIsFalse(t *testing.T) {
	// arrange
	statuses := []Status{
		ErrorDefault,
		Unknown,
		UpToDate,
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

func Test_CalculateGlobalStatus_EmptyList_ResultStatusIsUnknown(t *testing.T) {
	// arrange
	statuses := []Status{}

	// act
	result := CalculateGlobalStatus(statuses)

	// assert
	if result.Name() != Unknown.Name() {
		t.Fail()
		t.Logf("CalculateGlobalStatus returned %s instead of %s", result, Unknown)
	}
}

func Test_CalculateGlobalStatus_MixedStatuses_ResultStatusIsUnknown(t *testing.T) {
	// arrange
	statuses := []Status{
		ErrorDefault,
		Unknown,
		UpToDate,
	}

	// act
	result := CalculateGlobalStatus(statuses)

	// assert
	if result.Name() != Unknown.Name() {
		t.Fail()
		t.Logf("CalculateGlobalStatus returned %s instead of %s", result, Unknown)
	}
}

func Test_CalculateGlobalStatus_AllStatusesAreUpToDate_ResultStatusIsUpToDate(t *testing.T) {
	// arrange
	statuses := []Status{
		UpToDate,
		UpToDate,
		UpToDate,
	}

	// act
	result := CalculateGlobalStatus(statuses)

	// assert
	if result.Name() != UpToDate.Name() {
		t.Fail()
		t.Logf("CalculateGlobalStatus returned %s instead of %s", result, UpToDate)
	}
}

func Test_CalculateGlobalStatus_AllStatusesAreOutdated_ResultStatusIsOutdated(t *testing.T) {
	// arrange
	statuses := []Status{
		Outdated,
		Outdated,
		Outdated,
	}

	// act
	result := CalculateGlobalStatus(statuses)

	// assert
	if result.Name() != Outdated.Name() {
		t.Fail()
		t.Logf("CalculateGlobalStatus returned %s instead of %s", result, Outdated)
	}
}
