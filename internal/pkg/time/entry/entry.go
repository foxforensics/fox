package entry

import (
	"fmt"
	"strings"
	"time"
)

var replacer = strings.NewReplacer(
	`\`, `\\`, // mask backslashes
	`|`, `\|`, // mask pipes
	`/`, `\/`, // mask slashes
	`:`, `\:`, // mask colons
)

type Entry struct {
	Name    string    `json:"name,omitempty"`
	Inode   string    `json:"inode,omitempty"`
	Size    uint64    `json:"size"`
	Mode    string    `json:"mode,omitempty"`
	Mtime   time.Time `json:"m_time,omitempty"`
	Atime   time.Time `json:"a_time,omitempty"`
	Ctime   time.Time `json:"c_time,omitempty"`
	Btime   time.Time `json:"b_time,omitempty"`
	Anomaly bool      `json:"anomaly,omitempty"`
}

func (e Entry) String() string {
	return e.AsBody()
}

func (e Entry) AsBody() string {
	return fmt.Sprintf("0|%s|%s|%s|0|0|%d|%d|%d|%d|%d",
		replacer.Replace(e.Name),
		e.Inode,
		e.Mode,
		e.Size,
		e.Atime.UTC().Unix(),
		e.Mtime.UTC().Unix(),
		e.Ctime.UTC().Unix(),
		e.Btime.UTC().Unix(),
	)
}

/*
date,time,timezone,MACB,source,sourcetype,type,user,host,short,desc,version,filename,inode,notes,format,extra
07/04/2026,14:22:01,UTC,M...,FILE,NTFS $MFT,Content Modification Time,-,-,OS: /Windows/System32/config/SAM,NTFS $MFT entry, sequence: 3,2,/dev/sda2,1289,-,filestat,file_size: 262144; allocated: True
07/04/2026,14:22:05,UTC,....,REG,UserAssist,Last Run Time,jdoe,DESKTOP-01,UserAssist: chrome.exe,[UserAssist] Program execution: chrome.exe count: 42,2,/Users/jdoe/NTUSER.DAT,-,-,winreg_default,count: 42; application: chrome.exe
07/04/2026,14:25:33,UTC,..C.,EVT,WinEVTX,Change Time,SYSTEM,DESKTOP-01,[4624] Logon,An account was successfully logged on,2,/Windows/System32/winevt/Logs/Security.evtx,-,-,winevtx,event_identifier: 4624; source_name: Microsoft-Windows-Security-Auditing
*/

func (e Entry) AsCSV() string {
	return fmt.Sprintf("%s;%s;%d;%s;%s;%s;%s",
		replacer.Replace(e.Name),
		e.Inode,
		e.Size,
		e.Atime.UTC().Format(time.RFC3339),
		e.Mtime.UTC().Format(time.RFC3339),
		e.Btime.UTC().Format(time.RFC3339),
		e.Ctime.UTC().Format(time.RFC3339),
	)
}
