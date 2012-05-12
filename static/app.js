/* Javascript for adding and deleting rows.
 */

function add_row() {
  row = document.getElementById('table').insertRow(-1)
  row.id = "row" + next_row_index;
  row.innerHTML = ROW_TEMPLATE.replace(/%%/g, Number(next_row_index).toString());
  next_row_index++;
}

function delete_row(id) {
  child = document.getElementById(id);
  child.parentNode.removeChild(child);
}


var next_row_index = 1;
// %% is the placeholder
var ROW_TEMPLATE = '' +
'<td><input name="name%%" type="text">' +
'</td><td>' +
'<select name="weight%%">' +
'  <option value="5"> 5 </input>' +
'  <option value="10">10</input>' +
'  <option value="15">15</input>' +
'  <option value="20">20</input>' +
'  <option value="25">25</input>' +
'  <option value="30">30</input>' +
'  <option value="35">35</input>' +
'  <option value="40">40</input>' +
'  <option value="45">45</input>' +
'  <option value="50">50</input>' +
'  <option value="55">55</input>' +
'  <option value="60">60</input>' +
'  <option value="65">65</input>' +
'  <option value="70">70</input>' +
'  <option value="75">75</input>' +
'  <option value="80">80</input>' +
'  <option value="85">85</input>' +
'  <option value="90">90</input>' +
'  <option value="95">95</input>' +
'</td><td>' +
'<select name="arms%%">' +
'  <option value="0">0</input>' +
'  <option value="1">1</input>' +
'  <option value="2">2</input>' +
'  <option value="3">3</input>' +
'  <option value="4">4</input>' +
'  <option value="5">5</input>' +
'  <option value="6">6</input>' +
'  <option value="7">7</input>' +
'  <option value="8">8</input>' +
'  <option value="9">9</input>' +
'</select>' +
'</td><td>' +
'<select name="handles%%">' +
'  <option value="arms">Arms</input>' +
'  <option value="outer ground">Outer ground</input>' +
'  <option value="inner ground">Inner ground</input>' +
'  <option value="lat bar">Lat bar</input>' +
'</select>' +
'</td><td>' +
'<select name="handle_length%%">' +
'  <option value="short">Short</input>' +
'  <option value="long">Long</input>' +
"  <option value='doesnt_matter'>Doesn't matter</input>" +
'</select>' +
'</td><td>' +
'<select name="back%%">' +
'  <option value="flat" selected="selected">Flat</input>' +
'  <option value="curved">Curved</input>' +
"  <option value='doesnt_matter'>Doesn't matter</input>" +
'</select>' +
'</td><td>' +
'<select name="seat%%">' +
'  <option value="yes" selected="selected">Yes</input>' +
'  <option value="no">No</input>' +
"  <option value='doesnt_matter'>Doesn't matter</input>" +
'</select>' +
'</td><td>' +
'<a href="" onclick="delete_row(\'row%%\'); return false">' +
'  <img src="/static/red_x.png" /></a>' +
'</td>';
