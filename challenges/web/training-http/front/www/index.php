<?php
//load classes
require_once ("classes/TplBlock.php");
$tpl = new TplBlock();

require_once ("classes/LocaleManager.php");
$currentLocale = new LocaleManager();
$tpl->addVars(array("langForm" => $currentLocale->get_locale_form()));

$step = isset($_GET['step'])? $_GET['step'] : "0";

//fill the left menus links
$navMenus = array(
	array("href"	=> "/", "caption"	=> _('Home')),
	array("href"	=> "/?step=1", "caption"	=> _('HTTP Methods')),
	array("href"	=> "/?step=2", "caption"	=> _('HTTP Headers')),
    array("href"	=> "/?step=3", "caption"	=> _('HTTP Redirection')),
    array("href"	=> "/?step=4", "caption"	=> _('HTTP Cookies')),
);
foreach($navMenus as $menu){
	$tplNav = new TplBlock("nav");
	$tplNav->addVars($menu);
	$tpl->addSubBlock($tplNav);
}

switch ($step){

    case "1": //HTTP Methods
        $tpl->addVars( array("pageTitle" => htmlentities( _('HTTP Methods')) ) );

        //check if part of challenge success or not
        if (isset($_POST['login']) && isset($_POST['password'])) {
            $pmessage =  _('Well done!');
        } else {
            $pmessage ="<a href='/post.php'>" . _('This page') ."</a> " . _('expects a POST request with two parameters, login and password, but you won\'t find anywhere on the page to input these parameters.')
                     . "<br />" ._('Instead, you\'ll have to complete this request manually with curl (look at the -d option).') ."</p>";
        }

        //display the content

        $tplContent = new TplBlock();
        $content = $tplContent->addVars( array("pmessage"  => $pmessage) )
                    ->applyTplStr( LocaleManager::translate_tagged_parts(file_get_contents("templates/step1.html")) );
        $tpl->addVars( array("content" => $content) );

        break;

    case "2": // HTTP Headers
        $tpl->addVars( array("pageTitle" => htmlentities( _('HTTP Headers')) ) );

         //check if part of challenge success or not
         header('X-Pathwar-Header: ValidationCodeWillBeGivenAtANextStep'); 
         if ($_SERVER['HTTP_USER_AGENT'] == 'Sup3rH4x0RBr0ws3r') {
            $pmessage = _('W00t! Seems you are a real hacker!');
        }else {
            $pmessage = _('Try to change your User-Agent to <b>Sup3rH4x0RBr0ws3r</b> to get the passphrase !');
        }

        //display the content
        $tplContent = new TplBlock();
        $content = $tplContent->addVars( array("pmessage"  => $pmessage) )
                    ->applyTplStr( LocaleManager::translate_tagged_parts(file_get_contents("templates/step2.html")) );
        $tpl->addVars( array("content" => $content) );


        break;

    case "3":
        $tpl->addVars( array("pageTitle" => htmlentities( _('HTTP redirection')),
                             "content"    => LocaleManager::translate_tagged_parts( file_get_contents("templates/step3.html") )
        ) );
        break;

    case "4": //HTTP Cookies

        

        //check if part of challenge success or not
        if (!isset($_COOKIE['chocolate'])) {
            setcookie('chocolate', 'isbad');
        }

        if (isset($_COOKIE['chocolate']) && $_COOKIE['chocolate'] == 'isgood') {
            $pmessage = _('Congrats! The passphrase is') . ' __PASSPHRASE__';
        } else {
            $pmessage= _('Try to set the value of the cookie <b>chocolate</b> to <b>isgood</b> !');
        }

        $tplContent = new TplBlock();
        $content = $tplContent->addVars( array("pmessage"  => $pmessage) )
                    ->applyTplStr( LocaleManager::translate_tagged_parts(file_get_contents("templates/step4.html")) );

        $tpl->addVars( array("pageTitle"    => htmlentities( _('HTTP Cookies') )
                             ,"content"     =>  $content
                            )
                );

        break;
    case "redir":
        header('Location: /index.php?step=3');
        echo 'Congrats!';
        die();
        break;
    default: //HOME
        $tpl->addVars( array("pageTitle" => htmlentities( _('Training HTTP')),
                            "content"    => LocaleManager::translate_tagged_parts( file_get_contents("templates/home.html") )
        ));
        break;
}

echo $tpl->applyTplFile("templates/main.html");