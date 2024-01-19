function checkReload() {
  const xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
  let res = this.responseText;
  obj = JSON.parse(res);
  if (obj["refresh"]) {
  console.log("page has to be refeshed");
  window.location.reload()
  } else {
  console.log("no refesh needed");
  }
  console.log(res);
      console.log("reload answer");
}
  };
  xhttp.open("GET", "/_refresh", true);
  xhttp.send(); 	  
console.log("reload");
} 
window.setInterval(checkReload, 5000)
