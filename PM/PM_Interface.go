package PM

import (
	"time"
	"github.com/Esseh/retrievable"
	"github.com/Esseh/notorious-dev/AUTH"
	"github.com/Esseh/notorious-dev/CONTEXT"
)

func SendMessage(ctx CONTEXT.Context, email, title, content string){
	// Ensure Each User Exists
	ref := AUTH.EmailReference{}
	if ctx.User.Email == "" || retrievable.GetEntity(ctx, email , &ref) != nil {return}
	
	// Retrieve Each User's Header if it exists.
	SenderHeader := PrivateMessageHeader{}
	RecieverHeader := PrivateMessageHeader{}
	
	retrievable.GetEntity(ctx,int64(ctx.User.IntID),&SenderHeader)
	retrievable.GetEntity(ctx,ref.UserID,&SenderHeader)

	PrivateMSG := PrivateMessage{
		Sender:		ctx.User.Email,
		Receiver:	email,
		Title:		title,
		Content:	content,
		DateSent:	time.Now(),
	}
	
	// Send the Message
	key , _ := retrievable.PlaceEntity(ctx.Context,int64(0), &PrivateMSG)
	
	SenderHeader.Messages = append([]int64{key.IntID()},SenderHeader.Messages...) 
	RecieverHeader.Messages = append([]int64{key.IntID()},RecieverHeader.Messages...) 
	retrievable.PlaceEntity(ctx, int64(ctx.User.IntID), &SenderHeader)
	retrievable.PlaceEntity(ctx, ref.UserID, &RecieverHeader)
}