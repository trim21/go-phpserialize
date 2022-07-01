<?php

$val = array(
  'users' => array(
    0 => array(
      'id' => 1,
      'name' => 'sai',
    ),
    1 => array(
      'id' => 2,
      'name' => 'trim21',
    ),
  ),
  'obj' => array(
    'v' => 2,
    'e' => 3.14,
    'a long string name replace field name' => 'vvv',
  ),
);



$returnValue = serialize($val);

echo $returnValue;

// print_r(unserialize('a:2:{i:2;s:3:"two";i:1;s:3:"one";}'));
// echo $returnValue;
