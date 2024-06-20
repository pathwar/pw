<?php

//load classes
require_once ("classes/TplBlock.php");
$tpl = new TplBlock();

require_once ("classes/LocaleManager.php");
$currentLocale = new LocaleManager();
$tpl->addVars(
    array("langForm"    => $currentLocale->get_locale_form(),
          "source1"     => htmlentities(file_get_contents ("source_1.txt") ),
          "source2"     => highlight_string (file_get_contents ("source_2.txt"),true )
));

if(isset($_GET['pass'])){
    if($_GET['pass'] == md5("__RANDNUM0__")){
        $tplWin = new TplBlock("win");
        $tplWin ->addVars( array("PASSPHRASE0"  => '__PASSPHRASE0__') );
        $tpl->addSubBlock($tplWin);
    }else{
        $tplFail = new TplBlock("fail");
        $tpl->addSubBlock($tplFail);
    }
}

echo $tpl->applyTplStr( LocaleManager::translate_tagged_parts(file_get_contents("templates/main.html")) );