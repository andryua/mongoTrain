function funcSelect(sel) {
    var defgp = document.querySelectorAll(".default");
    var maingp = document.querySelectorAll(".main");
    var usrgp = document.querySelectorAll(".users");
    var keyarrdef = Object.keys(defgp);
    var keyarrmain = Object.keys(maingp);
    var keyarrusr = Object.keys(usrgp);
    switch (sel.value) {
        default:
            document.getElementById("collapsedit").style.display = "none";
            break;
        case "main":
            document.getElementById("collapsedit").style.display = "none";
            break;
        case "default":
            document.getElementById("collapsedit").style.display = "block";
            keyarrdef.forEach(function (key){
                defgp[key].style.display = "none";
            });
            keyarrmain.forEach(function (key){
                maingp[key].style.display = "block";
            });
            keyarrusr.forEach(function (key){
                usrgp[key].style.display = "none";
            });
            break;
        case "users":
            document.getElementById("collapsedit").style.display = "block";
            keyarrdef.forEach(function (key){
                defgp[key].style.display = "block";
            });
            keyarrmain.forEach(function (key){
                maingp[key].style.display = "none";
            });
            keyarrusr.forEach(function (key){
                usrgp[key].style.display = "none";
            });
            break;
    }
}
function funcType(type) {
    switch (type.class) {
        default:
            type.innerText = "none";
            break;
        case ".main":
            type.innerText = "головна";
            break;
        case ".default":
            type.innerText = "типова";
            break;
        case ".users":
            type.innerText = "користувацька";
            break;
    }
}