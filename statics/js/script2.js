function searchError() {
  var elemMsg = document.getElementById("error");
  elemMsg.innerHTML =
    '<div class="alert alert-warning alert-dismissible" id="alertfadeout" role="alert"><strong>WARNING</strong>: Search failed!</div>';
	 $(document).ready(function() {
    $('#alertfadeout').fadeIn(2000).delay(2000).fadeOut(2000);
  });
}

/*获取port和IP的排名*/
function getPortAndIPRank() {
$.post("/getPortAndIPRank","",function(data){
	portRanks=data.split("$")[0];
	ipRanks=data.split("$")[1];
	portRanka=portRanks.split("#");
	var phtml = '<div><table class="table table-striped"><thead><tr class="col-sm-2 col-md-2"><th class="col-sm-1 col-md-1">Port</th><th class="col-sm-1 col-md-1">Rank</th></tr></thead><tbody class="col-sm-2 col-md-2">';
	for (i = 0; i < portRanka.length; i++) {
	pr=portRanka[i].split(":");
	if(pr.length==2)
	{
		 phtml += '<tr ><td> <a style="text-decoration:none">' +  pr[0] + '</a></td><td>'+  pr[1] + '</td></tr>';
	}
    }
	document.getElementById("rankPort").innerHTML=phtml;

	ipRanka=ipRanks.split("#");
	var ihtml = '<div><table class="table table-striped"><thead><tr class="col-sm-2 col-md-2"><th class="col-sm-1 col-md-1">IP</th><th class="col-sm-1 col-md-1">Rank</th></tr></thead><tbody class="col-sm-2 col-md-2">';
	for (i = 0; i < ipRanka.length; i++) {
	ir=ipRanka[i].split(":");
	if(ir.length==2)
	{
		 ihtml += '<tr><td> <a style="text-decoration:none">' +  ir[0] + '</a></td><td>'+  ir[1] + '</td></tr>';
	}
    }
	document.getElementById("rankHost").innerHTML=ihtml;
	});
}
/*显示所有数据库*/
function getDatabase() {
 $.post("/getDatabase","",function(data){
     opetions=data.split("$");
	 var selectContent="";
	 for (i=0; i<opetions.length; i++ ){
	 selectContent=selectContent+"<option>"+opetions[i]+"</option>";
     }  
	 document.getElementById("ProjectSelect").innerHTML = selectContent;  
    });
}

/*清除显示结果*/
function clearLastResult() {
  document.getElementById("uploadBar").innerHTML = '';
    document.getElementById("status").innerHTML = '';
	 document.getElementById("NF_name").value = '';
}
/*文件上传*/
function file_upload() {
  var upfile = document.getElementById("lefile");
    document.getElementById("uploadBar").innerHTML = '';
    document.getElementById("status").innerHTML = '';

    var formData = new FormData();
    for (var i = 0; i < upfile.files.length; i++) {
      sendFile = "file" + i;
      formData.append(sendFile, upfile.files[i]);
    }
    var xmlhttp = new XMLHttpRequest();
    xmlhttp.upload.addEventListener("progress", progressHandler, false);
    xmlhttp.addEventListener("load", completeHandler, false);
    xmlhttp.addEventListener("error", errorHandler, false);
    xmlhttp.addEventListener("abort", abortHandler, false);
    xmlhttp.open("POST", "/upload", true);
    xmlhttp.send(formData);
  
}
function progressHandler(event) {
  var percent = (event.loaded / event.total) * 100;
  document.getElementById("uploadBar").innerHTML = '<h4>Upload ...</h4><div class="progress"><div class="progress-bar progress-bar-striped active" role="progressbar" style="width: ' + Math.round(percent) + '%;">' + Math.round(percent) + '%</div></div>';
}

var parse_status = false;

function completeHandler(event) {
  if (event.target.responseText == "FAIL") {
    document.getElementById("status").innerHTML = '<div class="alert alert-danger"><strong>ERROR</strong>: Upload Failed!</div>';
  }
  if (event.target.responseText == "SUCCESS") {
    parse_status = false
    document.getElementById("uploadBar").innerHTML = '<h4>Upload ...</h4><div class="progress"><div class="progress-bar progress-bar-success progress-bar-striped" role="progressbar" style="width: 100%;">Waiting ...</div></div>';
	setTimeout(loop, 2000);
    var loop = function() {
      if (parse_status == false) {
        setTimeout(loop, 2000);
      }
      parseEVTX();
    }
    loop();
  }
}

function errorHandler(event) {
  document.getElementById("status").innerHTML = '<div class="alert alert-danger"><strong>ERROR</strong>: Upload Failed!</div>';
}

function abortHandler(event) {
  document.getElementById("status").innerHTML = '<div class="alert alert-info">Upload Aborted</div>';
}

function parseEVTX() {
  var xmlhttp2 = new XMLHttpRequest();
  xmlhttp2.open("GET", "/status");
  xmlhttp2.send();
  xmlhttp2.onreadystatechange = function() {
    if (xmlhttp2.readyState == 4) {
      if (xmlhttp2.status == 200) {
        var statusdata = xmlhttp2.responseText;
      if (statusdata=="OK") {
          document.getElementById("uploadBar").innerHTML = '<h4>Parsing  ...</h4><div class="progress"><div class="progress-bar progress-bar-success progress-bar-striped" role="progressbar" style="width: 100%;">SUCCESS</div></div>';
          document.getElementById("status").innerHTML = '<div class="alert alert-info"><strong>Import Success</strong>: You need to reload the web page.</div>';
          parse_status = true;
		  getPortAndIPRank();
        } 
      }
    }
  }
}


