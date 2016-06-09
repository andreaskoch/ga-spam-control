package spamcontrol

type trainer interface {
	// GetTrainingData returns a set of training data for the given accountID.
	// Returns an error if the training data could not be fetched.
	GetTrainingData(accountID string, numberOfDays int) (MachineLearningModel, error)
}

// MachineLearningModelTrainer returns pre-rated training data for
// for training machine learning models.
type MachineLearningModelTrainer struct {
	analyticsDataProvider analyticsDataProvider
	spamDomainRepository  SpamDomainRepository
}

// GetTrainingData returns a set of training data for the given accountID.
// Returns an error if the training data could not be fetched.
func (trainer *MachineLearningModelTrainer) GetTrainingData(accountID string, numberOfDays int) (MachineLearningModel, error) {
	analyticsData, analyticsDataError := trainer.analyticsDataProvider.GetAnalyticsData(accountID, numberOfDays)
	if analyticsDataError != nil {
		return MachineLearningModel{}, analyticsDataError
	}

	// convert the analytics data to a machine learning model
	machineLearningModel := analyticsDataToMachineLearningModel(analyticsData)

	// remove duplicates
	machineLearningModel.Rows = removeDuplicatesFromTable(machineLearningModel.Rows)

	// add spam rating
	spamDomainNames, spamDomainsError := trainer.spamDomainRepository.GetSpamDomains()
	if spamDomainsError != nil {
		return MachineLearningModel{}, spamDomainsError
	}

	machineLearningModel.ColumnNames = append(machineLearningModel.ColumnNames, "isSpam")
	for index, row := range machineLearningModel.Rows {

		// Check if the row is spam or not.
		// Note, this will only work with known
		// referrer spam domain names. A manual review is advised.
		isSpam := trainingdataFalse
		for _, spamDomainName := range spamDomainNames {

			// check if the domain name matches a known referrer spam domain
			domainName := row[len(row)-1]
			if domainName == spamDomainName {
				isSpam = trainingdataTrue
				break
			}

			// no match
		}

		// append the spam rating
		machineLearningModel.Rows[index] = append(machineLearningModel.Rows[index], isSpam)

	}

	return MachineLearningModel(machineLearningModel), nil
}
