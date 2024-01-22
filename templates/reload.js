function checkReload() {
  const xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
      let obj = JSON.parse(this.responseText);
      if (obj["refresh"]) {
        console.log("page has to be refeshed");
        window.location.reload()
      }
    }
  };
  xhttp.open("GET", "/_refresh", true);
  xhttp.send(); 	  
} 

window.setInterval(checkReload, 1000)
