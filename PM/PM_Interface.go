package PM

import (
	"time"
	"github.com/Esseh/retrievable"
	"github.com/Esseh/notorious-dev/AUTH"
	"github.com/Esseh/notorious-dev/CONTEXT"
	"github.com/Esseh/notorious-dev/NOTIFICATION"
)

func Notify(ctx CONTEXT.Context,email string){
	er := AUTH.EmailReference{}
	err := retrievable.GetEntity(ctx,email,&er)
	if err != nil { return }
	n := NOTIFICATION.Notifications{}
	retrievable.GetEntity(ctx,er.UserID,&n)
	n.NotificationsPending += 1
	n.Notifications = append([]string{ctx.User.Email + " sent you a private message"},n.Notifications...)
	retrievable.PlaceEntity(ctx,er.UserID,&n)
}

func SendMessage(ctx CONTEXT.Context, email, title, content string){
	// Ensure Each User Exists
	ref := AUTH.EmailReference{}
	if ctx.User.Email == "" || retrievable.GetEntity(ctx, email , &ref) != nil {return}
	
	// Retrieve Each User's Header if it exists.
	SenderHeader := PrivateMessageHeader{}
	RecieverHeader := PrivateMessageHeader{}
	
	retrievable.GetEntity(ctx,int64(ctx.User.IntID),&SenderHeader)
	retrievable.GetEntity(ctx,ref.UserID,&RecieverHeader)

	PrivateMSG := PrivateMessage{
		Sender:		ctx.User.Email,
		Receiver:	email,
		Title:		title,
		Content:	content,
		DateSent:	time.Now(),
	}
	
	// Send the Message
	key , _ := retrievable.PlaceEntity(ctx.Context,int64(0), &PrivateMSG)
	
	RecieverHeader.Messages = append([]int64{key.IntID()},RecieverHeader.Messages...)
	SenderHeader.Messages = append([]int64{key.IntID()},SenderHeader.Messages...)
	retrievable.PlaceEntity(ctx, int64(ctx.User.IntID), &SenderHeader)
	retrievable.PlaceEntity(ctx, ref.UserID, &RecieverHeader)
}

func RetrieveMessages(ctx CONTEXT.Context,pageWidth,pageNumber int)[]PrivateMessage{
	header := PrivateMessageHeader{}
	retrievable.GetEntity(ctx,int64(ctx.User.IntID),&header)
	lowerBound := pageWidth*pageNumber
	upperBound := pageWidth*(pageNumber+1)
	if lowerBound >= len(header.Messages){ return []PrivateMessage{} }
	if upperBound >  len(header.Messages){ upperBound = len(header.Messages) }
	gatheredMessages := header.Messages[lowerBound:upperBound]
	output := []PrivateMessage{}
	for _,v := range gatheredMessages {
		message := PrivateMessage{}
		if retrievable.GetEntity(ctx,v,&message) != nil { continue }
		output = append(output,message)
	}
	return output
}

func GetPageNumbers(ctx CONTEXT.Context,pageWidth,currentPage int) []int64 {
	header := PrivateMessageHeader{}
	retrievable.GetEntity(ctx,int64(ctx.User.IntID),&header)
	remainingNumbers := 4
	t := len(header.Messages) - 1;
	if t < remainingNumbers { remainingNumbers = t }
	if t < 0 { return []int64{} }
	iterator := 1
	output := []int64{int64(currentPage)}
	for remainingNumbers != 0 {
		if iterator > 20 { break }
		if currentPage - iterator > 0 {
			// Prepend Page Number
			output = append([]int64{int64(currentPage - iterator)},output...)		
			remainingNumbers--
		}
		if currentPage + iterator < len(header.Messages) {
			// Append Page Number
			output = append(output, int64(currentPage + iterator))			
			remainingNumbers--
		}
		iterator += 1
	}
	return output
}