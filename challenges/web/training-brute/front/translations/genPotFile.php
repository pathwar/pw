<?php

$dirs = array("../www", "../www/templates");
$potContent = "";
$allLines = array();
foreach( $dirs as $dir)
{
    $files = scandir($dir);
    foreach($files as $file) {

      if(is_file( $dir. "/" . $file )){
          
        preg_match_all("/_\(\'(.*?)\'\)/is", file_get_contents($dir. "/" . $file ), $m);
        foreach($m[1] as $stringToTranslate){
            if(!in_array($stringToTranslate, $allLines) ){ //no duplicated lines
                $allLines[] = $stringToTranslate;

                if(strstr($stringToTranslate, PHP_EOL)) {
                    //there is a carriage return
                    $potContent.='msgid ""' . "\n";
                    $lines = explode(PHP_EOL,$stringToTranslate);
                    foreach($lines as $line){
                        $potContent.= '"' . str_replace('"','\\"',$line ). '"' .PHP_EOL;
                    }
                    $potContent.='msgstr ""' . PHP_EOL. PHP_EOL;
                }else{
                    $potContent.='msgid "' . str_replace('"','\\"',$stringToTranslate ) .'"' . "\n";
                    $potContent.='msgstr ""' . PHP_EOL. PHP_EOL;
                }

            }

        }


      }
    }
}

file_put_contents("messages.pot", $potContent);