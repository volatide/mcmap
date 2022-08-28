package database 

import (
    "fmt"

    "github.com/go-pg/pg/v10"
)


/*
// PRISMA SCHEMA:
model Server {
  id      	String      @id @default(cuid())
  nick    	String?
  host    	String
  port    	Int
  disabled	Boolean		@default(false)
  entries 	ScanEntry[]
}

model ScanEntry {
  id            String   @id @default(cuid())
  server        Server   @relation(fields: [serverId], references: [id])
  serverId      String
  date          DateTime
  motd          String
  playerCount   Int
  maxPlayers    Int
  version       String
  protocol      Int
  icon          String
  playerSamples ScanEntryPlayerRelation[]
}

model ScanEntryPlayerRelation {
	player		Player @relation(fields: [playerId], references: [id])
	playerId	String	
	scanEntry	ScanEntry @relation(fields: [scanEntryId], references: [id])
	scanEntryId	String	

	@@id([playerId, scanEntryId])
}

model Player {
  id          String     @id @default(cuid())
  name        String
  uuid        String
  scanEntries ScanEntryPlayerRelation[] 
}
*/

type Server struct {
	id			string
	nick		string
	host		string
	port		int
	disabled	bool
	entries		[]ScanEntry
}

type Player struct {
	id			string
	uuid		string
	name		string
	scanEntries []ScanEntryPlayerRelation
}

type ScanEntryPlayerRelation struct {
	player		Player
	playerId	string

	scanEntry	ScanEntry
	scanEntryId	string
}

type ScanEntry struct {
    id				string	
	server			Server	
	serverId		string
	date			uint32
	motd			string
	playerCount		int	
	maxPlayers		int
	version			string
	protocol		int
	icon			string
	playerSamples	[]ScanEntryPlayerRelation
}

func (entry ScanEntry) String() string {
	return fmt.Sprintf("[%s] id=%s motd=%s version=%s players=(%d/%d)", entry.server.host, entry.id, entry.motd, entry.version, entry.playerCount, entry.maxPlayers)
}
