package globals

import (
	"strconv"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
)

var HppServerAccount *nex.Account

func AccountDetailsByPID(pid types.PID) (*nex.Account, *nex.Error) {
	if pid.Equals(HppServerAccount.PID) {
		return HppServerAccount, nil
	}

	password, errorCode := PasswordFromPID(pid)
	if errorCode != 0 {
		return nil, nex.NewError(errorCode, "Failed to get password from PID")
	}

	account := nex.NewAccount(pid, strconv.Itoa(int(pid)), password)

	return account, nil
}

func AccountDetailsByUsername(username string) (*nex.Account, *nex.Error) {
	if username == HppServerAccount.Username {
		return HppServerAccount, nil
	}

	pidInt, err := strconv.Atoi(username)
	if err != nil {
		return nil, nex.NewError(nex.ResultCodes.RendezVous.InvalidUsername, "Invalid username")
	}

	pid := types.NewPID(uint64(pidInt))

	password, errorCode := PasswordFromPID(pid)
	if errorCode != 0 {
		return nil, nex.NewError(errorCode, "Failed to get password from PID")
	}

	account := nex.NewAccount(pid, username, password)

	return account, nil
}
