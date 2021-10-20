<?php
$result = print_r($_SERVER, true);
file_put_contents("result.log", print_r($_SERVER, true));
echo $result;
echo file_get_contents("php://input");
echo print_r($_GET, true);
echo print_r($_POST, true);
