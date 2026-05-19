package client

import (
	"errors"
	"fmt"
)

var errClientOnly = errors.New("client-only command")

type commandHelp struct {
	usage       string
	protocol    string
	description string
}

var commandHelpList = []commandHelp{
	{usage: "/help", protocol: "(client)", description: "Show this help"},
	{usage: "/nick <name>", protocol: "LOGIN", description: "Set your nickname"},
	{usage: "/signup <name>", protocol: "SIGN_UP", description: "Sign up (server stub)"},
	{usage: "/quit", protocol: "LOGOUT", description: "Disconnect"},
	{usage: "/create <room>", protocol: "CREATE_ROOM", description: "Create a room"},
	{usage: "/join <room>", protocol: "JOIN_ROOM", description: "Join a room"},
	{usage: "/leave [room]", protocol: "LEAVE_ROOM", description: "Leave a room"},
	{usage: "/deleteroom <room>", protocol: "DELETE_ROOM", description: "Delete a room you own"},
	{usage: "/editroom <room> <new> [max]", protocol: "EDIT_ROOM", description: "Rename room and/or set max size"},
	{usage: "/members <room>", protocol: "GET_ROOM_MEMBERS", description: "List members in a room"},
	{usage: "/kick <room> <user>", protocol: "DELETE_MEMBER", description: "Remove a member (owner)"},
	{usage: "/invite <room> <user>", protocol: "SEND_INVITE_REQUEST", description: "Invite user to a room"},
	{usage: "/invites [room]", protocol: "SEE_GROUP_INVITE_REQUEST", description: "List pending room invites"},
	{usage: "/acceptinvite <room>", protocol: "ACCEPT_GROUP_INVITE_REQUEST", description: "Accept a room invite"},
	{usage: "/declineinvite <room>", protocol: "DELETE_GROUP_INVITE_REQUEST", description: "Decline a room invite"},
	{usage: "/rooms", protocol: "LIST_ROOMS", description: "List public rooms"},
	{usage: "/myrooms", protocol: "LIST_MY_ROOMS", description: "List rooms you joined"},
	{usage: "/msg <text>", protocol: "SEND_MSG", description: "Send message to current room"},
	{usage: "<text>", protocol: "SEND_MSG", description: "Same as /msg when in a room"},
	{usage: "/sendfile <path>", protocol: "SEND_FILE", description: "Send a file (server stub)"},
	{usage: "/friendadd <user>", protocol: "SEND_FRIEND_REQUEST", description: "Send friend request"},
	{usage: "/friendaccept <user>", protocol: "ACCEPT_FRIEND_REQUEST", description: "Accept friend request"},
	{usage: "/friendrequests", protocol: "SEE_FRIEND_REQUEST", description: "List pending friend requests"},
	{usage: "/friends", protocol: "GET_FRIENDS", description: "List friends"},
	{usage: "/friendremove <user>", protocol: "DELETE_FRIEND", description: "Remove a friend"},
	{usage: "/dm <user> <msg>", protocol: "MESSAGE_FRIEND", description: "Direct message a friend"},
	{usage: "/block <user>", protocol: "BLOCK_USER", description: "Block a user"},
	{usage: "/unblock <user>", protocol: "UNBLOCK_USER", description: "Unblock a user"},
	{usage: "/status <user>", protocol: "GET_USER_STATUS", description: "Get user online status"},
}

func printClientHelp() {
	fmt.Printf("\nCommands (client → protocol action):\n")
	for _, c := range commandHelpList {
		fmt.Printf("  %-28s %-26s %s\n", c.usage, c.protocol, c.description)
	}
	fmt.Printf("\nServer prompts (no slash command):\n")
	fmt.Printf("  %-28s %-26s %s\n", "(prompt)", "SET_ROOM_PASSWORD", "Set password when creating a room")
	fmt.Printf("  %-28s %-26s %s\n", "(prompt)", "GET_ROOM_PASSWORD", "Enter password to join a private room")
	fmt.Printf("\n")
}

func printBootstrap() {
	fmt.Printf("\n╔════════════════════════════════════════╗\n")
	fmt.Printf("║         CLI Chat Client                ║\n")
	fmt.Printf("╠════════════════════════════════════════╣\n")
	fmt.Printf("║  Type /help for all commands           ║\n")
	fmt.Printf("║  /nick <name>  — login                 ║\n")
	fmt.Printf("║  /join <room>  — join a room           ║\n")
	fmt.Printf("║  plain text    — chat in current room  ║\n")
	fmt.Printf("╚════════════════════════════════════════╝\n")
}
