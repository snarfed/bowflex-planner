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
'</td><td><input name="weight%%" type="text" size="3">' +
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
'  <option value="Outer ground">Outer ground</input>' +
'  <option value="Inner ground">Inner ground</input>' +
'  <option value="Lat bar">Lat bar</input>' +
'</select>' +
'</td><td>' +
'<select name="handle_length%%">' +
'  <option value="Short">Short</input>' +
'  <option value="Long">Long</input>' +
'</select>' +
'</td><td>' +
'<select name="back%%">' +
'  <option value=""></input>' +
'  <option value="Flat">Flat</input>' +
'  <option value="Curved">Curved</input>' +
'</select>' +
'</td><td>' +
'<input type="checkbox" name="seat%%" value="yes" /> ' +
'</td><td>' +
'<a href="" onclick="delete_row(\'row%%\'); return false">' +
'  <img src="/static/red_x.png" /></a>' +
'</td>';
