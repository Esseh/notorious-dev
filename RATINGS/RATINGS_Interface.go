package RATINGS

import (
	"strconv"
	"github.com/Esseh/retrievable"
	"github.com/Esseh/notorious-dev/USERS"
	"github.com/Esseh/notorious-dev/NOTES"
	"github.com/Esseh/notorious-dev/CONTEXT"
)

func GetRating(ctx CONTEXT.Context) string {
	noteID, err := strconv.ParseInt(ctx.Req.FormValue("NoteID"),10,64)
	if err != nil { return `{"success":false,"code":0}` }
	if retrievable.GetEntity(ctx,noteID,&NOTES.Note{}) != nil { return `{"success":false,"code":1}`}
	r := Rating{}
	if retrievable.GetEntity(ctx,noteID,&r) != nil {
		retrievable.PlaceEntity(ctx,noteID,&r)
	}
	if len(r.Entries) == 0{
		return `{"success":true,"totalRating":0,"code":-1}`
	}
	var totalValue float64
	for _,v := range r.Entries {
		totalValue += float64(v.Score)
	}
	totalValue /= float64(len(r.Entries))
	outputValue := strconv.FormatFloat(totalValue,'E',-1,64)
	if outputValue[1] == 'E' { outputValue = outputValue[:1] 
	} else { outputValue = outputValue[:3] }
	return `{"success":true,"totalRating":`+outputValue+`,"code":-1}`
}

func SetRating(ctx CONTEXT.Context) string{
	r := Rating{}
	noteID, err := strconv.ParseInt(ctx.Req.FormValue("NoteID"),10,64)
	if err != nil { return `{"success":false,"code":0}` }
	value, err := strconv.ParseInt(ctx.Req.FormValue("RatingValue"),10,64)
	if err != nil { return `{"success":false,"code":0}` }
	if (*(ctx.User) == USERS.User{}) { return `{"success":false,"code":0}` }
	if retrievable.GetEntity(ctx,noteID,&r) != nil { return `{"success":false,"code":1}`}
	if value < 1 || value > 5 { return `{"success":false,"code":0}` }
	IsNewEntry := true
	for i,v := range r.Entries {
		if retrievable.IntID(v.UserID) == ctx.User.IntID {
			r.Entries[i].Score = value
			IsNewEntry = false
			break
		}
	}
	if IsNewEntry { r.Entries = append(r.Entries,Pair{int64(ctx.User.IntID),value}) }
	if _ , placeErr := retrievable.PlaceEntity(ctx,noteID,&r); placeErr != nil { return `{"success":false,"code":1}` }
	return `{"success":true,"code":-1}`
}