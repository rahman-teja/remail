package remail

import (
	"fmt"
)

var (
	ErrEmailNotFound                             error = fmt.Errorf("email: Not found")
	ErrSenderIsRequired                          error = fmt.Errorf("email: Sender is required")
	ErrDestinationIsRequired                     error = fmt.Errorf("email: Destination is required")
	ErrConfigNotFound                            error = fmt.Errorf("email: Config found")
	ErrCodeMessageRejected                       error = fmt.Errorf("email: Messasge Rejected")
	ErrCodeMailFromDomainNotVerifiedException    error = fmt.Errorf("email: Mail From Domain Not Verified")
	ErrCodeConfigurationSetDoesNotExistException error = fmt.Errorf("email: Configuration Set Does Not Exist")
)
