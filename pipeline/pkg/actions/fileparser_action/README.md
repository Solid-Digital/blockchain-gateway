# File parser action

## Input
The input map can contain the following values:


<table>
    <tr>
        <td><b>Key</b></td>
        <td><b>Description</b></td>
        <td><b>Allowed values</b></td>
    </tr>
    <tr>
        <td>Filetype</td>
        <td>The type of file to parse</td>
        <td>"csv"</td>
    </tr>
    <tr>
		<td>File</td>
		<td>Raw file as byte array ([]byte)</td>
		<td>Any []byte</td>
	</tr>
	<tr>
		<td>Header</td>
		<td>Indicate whether a CSV file contains a header. If it does the result will contain the keys.</td>
		<td>true or false (optional)</td>
	</tr>   
</table>

## Output

* Messages:
The output of the file parser action is a map with messages. 
Messages contains an array of map[string]interface{}.

### CSV

In case the CSV does not contain headers the result will be put in a map based on the order of the columns. The keys in the map 
will be the column number prefixed with 'col-', e.g. the first column of a csv file without headers will be part of a message output map under the key 'col-1'.  