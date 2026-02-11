package rules

import (
	"slices"
	"strings"

	"github.com/bradleyjkemp/sigma-go"
)

var Default = []byte(`
title: Critical Events
logsource:
  product: fox
detection: 
  selection:
    - PRIORITY:
      - 0 # System is unusable
      - 1 # Action must be taken immediately
      - 2 # Critical conditions
      - 3 # Error conditions
    - EventID:
      - 1102 # The audit log was cleared
      - 4624 # An account was successfully logged on
      - 4625 # An account failed to log on
      - 4648 # A logon was attempted using explicit credentials
      - 4663 # An attempt was made to access an object
      - 4664 # An attempt was made to create a hard link
      - 4670 # Permissions on an object were changed
      - 4672 # Special privileges assigned to new logon
      - 4673 # A privileged service was called
      - 4674 # An operation was attempted on a privileged object
      - 4688 # A new process has been created
      - 4690 # An attempt was made to duplicate a handle to an object
      - 4692 # Backup of data protection master key
      - 4696 # A primary token was assigned to process
      - 4697 # A service was installed in the system
      - 4715 # The audit policy on an object was changed
      - 4717 # System security access was granted
      - 4718 # System security access was revoked
      - 4719 # System audit policy was changed
      - 4720 # A user account was created
      - 4726 # A user account was deleted
      - 4732 # A member was added to a security group
      - 4733 # A member was removed from a security group
      - 4738 # A user account was modified
      - 4739 # Domain Policy was changed
      - 4740 # A user account was locked out
      - 4756 # A member was added to a privileged group
      - 4757 # A member was removed from a privileged group
      - 4768 # A Kerberos authentication ticket was requested
      - 4769 # A Kerberos service ticket was requested
      - 4771 # Kerberos preauthentication failed
      - 4776 # The computer attempted to validate credentials
      - 4778 # A session was reconnected to a Window Station
      - 4779 # A session was disconnected from a Window Station
      - 4902 # The per user audit policy table was created
      - 4904 # An attempt was made to register a security event source
      - 4905 # An attempt was made to unregister a security event source
      - 5140 # A network share object was accessed
      - 5145 # A network share object was checked
      - 5148 # The Windows Filtering Platform has detected a DoS attack      
      - 5156 # The Windows Filtering Platform has permitted a connection
      - 5157 # The Windows Filtering Platform has blocked a connection
      - 5158 # The Windows Filtering Platform has permitted a bind to a local port
      - 6272 # Network Policy Server granted access
      - 6273 # Network Policy Server denied access
      - 7036 # The service entered the stopped/running state
      - 7045 # A service was installed in the system
  condition: selection
`)

var supported = []string{
	"fox",
	"linux",
	"windows",
}

func IsSupported(r *sigma.Rule) bool {
	return slices.Contains(supported, strings.ToLower(r.Logsource.Product))
}
