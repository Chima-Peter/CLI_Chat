package client

import (
	"errors"
	"fmt"
)

var errClientOnly = errors.New("client-only command")

type commandHelp struct {
	group       string
	usage       string
	protocol    string
	description string
}

var commandHelpList = []commandHelp{
	{group: "General", usage: "/help", protocol: "(client)", description: "Show this help"},

	{group: "Auth", usage: "/auth/login <name>", protocol: "LOGIN", description: "Set your nickname"},
	{group: "Auth", usage: "/auth/signup <name>", protocol: "SIGN_UP", description: "Sign up (server stub)"},
	{group: "Auth", usage: "/auth/logout", protocol: "LOGOUT", description: "Disconnect"},

	{group: "Room", usage: "/room/create <name>", protocol: "CREATE_ROOM", description: "Create a room"},
	{group: "Room", usage: "/room/join <name>", protocol: "JOIN_ROOM", description: "Join a room"},
	{group: "Room", usage: "/room/leave [name]", protocol: "LEAVE_ROOM", description: "Leave a room"},
	{group: "Room", usage: "/room/delete <name>", protocol: "DELETE_ROOM", description: "Delete a room you own"},
	{group: "Room", usage: "/room/edit <room> <new> [max]", protocol: "EDIT_ROOM", description: "Rename room and/or set max size"},
	{group: "Room", usage: "/room/members <room>", protocol: "GET_ROOM_MEMBERS", description: "List members in a room"},
	{group: "Room", usage: "/room/kick <room> <user>", protocol: "DELETE_MEMBER", description: "Remove a member (owner)"},
	{group: "Room", usage: "/room/invite <room> <user>", protocol: "SEND_INVITE_REQUEST", description: "Invite a user to a room"},
	{group: "Room", usage: "/room/invites [room]", protocol: "SEE_GROUP_INVITE_REQUEST", description: "List users you invited (owner)"},
	{group: "Room", usage: "/room/list", protocol: "LIST_ROOMS", description: "List public rooms"},
	{group: "Room", usage: "/room/mine", protocol: "LIST_MY_ROOMS", description: "List rooms you joined"},

	{group: "Invite", usage: "/invite/mine", protocol: "LIST_MY_ROOM_INVITES", description: "List rooms that invited you"},
	{group: "Invite", usage: "/invite/accept <room>", protocol: "ACCEPT_GROUP_INVITE_REQUEST", description: "Accept a room invite"},
	{group: "Invite", usage: "/invite/decline <room>", protocol: "DELETE_GROUP_INVITE_REQUEST", description: "Decline a room invite"},

	{group: "Chat", usage: "/chat/send <message>", protocol: "SEND_MSG", description: "Send message to current room"},
	{group: "Chat", usage: "<text>", protocol: "SEND_MSG", description: "Same as /chat/send when in a room"},
	{group: "Chat", usage: "/chat/file <path>", protocol: "SEND_FILE", description: "Send a file (server stub)"},
	{group: "Chat", usage: "/chat/dm <user> <msg>", protocol: "MESSAGE_FRIEND", description: "Direct message a friend"},

	{group: "Friend", usage: "/friend/add <user>", protocol: "SEND_FRIEND_REQUEST", description: "Send friend request"},
	{group: "Friend", usage: "/friend/accept <user>", protocol: "ACCEPT_FRIEND_REQUEST", description: "Accept friend request"},
	{group: "Friend", usage: "/friend/remove <user>", protocol: "DELETE_FRIEND", description: "Remove a friend"},
	{group: "Friend", usage: "/friend/requests", protocol: "SEE_FRIEND_REQUEST", description: "List incoming friend requests"},
	{group: "Friend", usage: "/friend/list", protocol: "GET_FRIENDS", description: "List friends"},

	{group: "User", usage: "/user/block <user>", protocol: "BLOCK_USER", description: "Block a user"},
	{group: "User", usage: "/user/unblock <user>", protocol: "UNBLOCK_USER", description: "Unblock a user"},
	{group: "User", usage: "/user/status <user>", protocol: "GET_USER_STATUS", description: "Get user online status"},
}

func printClientHelp() {
	fmt.Printf("\nCommands (grouped like API endpoints):\n")
	lastGroup := ""
	for _, c := range commandHelpList {
		if c.group != lastGroup {
			fmt.Printf("\n  [%s]\n", c.group)
			lastGroup = c.group
		}
		fmt.Printf("    %-32s %-26s %s\n", c.usage, c.protocol, c.description)
	}
	fmt.Printf("\nServer prompts (no slash command):\n")
	fmt.Printf("  %-32s %-26s %s\n", "(prompt)", "SET_ROOM_PASSWORD", "Set password when creating a room")
	fmt.Printf("  %-32s %-26s %s\n", "(prompt)", "GET_ROOM_PASSWORD", "Enter password to join a private room")
	fmt.Printf("\n")
}

func printBootstrap() {
	fmt.Printf("\n╔════════════════════════════════════════╗\n")
	fmt.Printf("║         CLI Chat Client                ║\n")
	fmt.Printf("╠════════════════════════════════════════╣\n")
	fmt.Printf("║  Type /help for all commands           ║\n")
	fmt.Printf("║  /auth/login <name>                    ║\n")
	fmt.Printf("║  /room/join <room>                     ║\n")
	fmt.Printf("║  plain text — chat in current room     ║\n")
	fmt.Printf("╚════════════════════════════════════════╝\n")
}
