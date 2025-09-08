package pkg

import (
	"fmt"
	"strconv"
	"strings"
)


func Id_to_file_name(id int64) string {
    name := fmt.Sprintf("%06d.data", id)
    return name
}
func Id_to_hint_name(id int64) string {
    name := fmt.Sprintf("%06d.hint", id)
    return name
}


func File_name_to_Id(name string) (int64, error) {
	base := strings.TrimSuffix(name , ".data");
	id , err := strconv.ParseInt(base , 10 , 64);
	if(err != nil){
		return 0, fmt.Errorf("error while generating id : %w",err);
	}
	return id , nil;
}
