<?php

/*
* This file is for PATHWAR challanges translation
* The code here won't help you to solve the challenge.
* If you find a flaw to exploit here... We didn't do it on purpose.
*/

class LocaleManager
{   
    const LOCALES_DOMAIN = "messages";
    const LOCALES_PATH = "../translations"; //relative path from web server root
    private $currentLocale = "en_US";
    private $availablesLocales = array();
    
    public function __construct()
    {
        $this->init_available_locales();

        if ( isset( $_POST['locale'] ) ){
            setcookie("locale", $_POST['locale'], strtotime( '+30 days' ) , "/", $_SERVER['SERVER_NAME'], 1);
            $this->set_locale( $_POST['locale'] );
            header("Refresh:0");
            die();
        }
        $this->init_current_locale();

    }
    private function get_available_locales()
    {
        if( empty( $this->availablesLocales ) ){
            $this->init_available_locales();
        }
        return $this->availablesLocales;
    }
    private function init_available_locales()
    {
        $this->availablesLocales = array();
        $cdir = scandir( self::LOCALES_PATH );
        
        foreach( $cdir as $directory ){
            if(in_array( $directory, ResourceBundle::getLocales('') )){
                $this->availablesLocales[] = $directory;
            }
        }
    }

    private function init_current_locale()
    {
        /*
        * priority order to choose the language: Cookie, navigator first language, default (en_US)
        */
        
        if(isset( $_COOKIE['locale'] ) && in_array( $_COOKIE['locale'], $this->availablesLocales) ){
            $this->set_locale( $_COOKIE['locale'] );
            return;
        }

        $browserLocale = substr($_SERVER['HTTP_ACCEPT_LANGUAGE'], 3, 5);
        if( in_array($browserLocale, $this->availablesLocales )){
            $this->set_locale( $browserLocale );
            return;
        }

        //not exactly the good locale, but maybe have the good language:
        $browserLocaleShort = substr($_SERVER['HTTP_ACCEPT_LANGUAGE'], 0, 2);
        foreach($this->availablesLocales as $availableLocale){
            if( $browserLocaleShort == substr($availableLocale,0,2) ){
                $this->set_locale( $availableLocale );
                return;
            }
        }
        $this->set_locale("en_US");
        return;

    }

    private function set_locale($locale)
    {
        $this->currentLocale = $locale;
        putenv('LC_ALL=' . $this->currentLocale.".UTF-8");
        setlocale(LC_ALL, $this->currentLocale.".UTF-8");
        bindtextdomain(self::LOCALES_DOMAIN, self::LOCALES_PATH);
        textdomain(self::LOCALES_DOMAIN);
    }

    public function get_locale_form()
    {
        $form = '<form method="POST"><select name="locale" onchange="this.form.submit();">';
        foreach( $this-> get_available_locales() as $availableLocale )
        {
            
            $form.='<option value="' . htmlentities($availableLocale) .'" '.
                    (($availableLocale == $this->currentLocale) ? 'selected="selected"' : ''). ">".
                    htmlentities($availableLocale) ."</option>";
        }
        $form.='</select></form>';
        return $form;
    }
    public function get_langShort(){

        return substr( $this->currentLocale, 0, 2 );
    }

    /*
    * extract tagged parts with _(' Content to translate ')
    * and pass them on the gettext function
    */
    public static function translate_tagged_parts (string $text)
    {
        $regex = "/_\(\'(.*?)\'\)/is";
        return preg_replace_callback(
            $regex,
            function ($m) {
                return _($m[1]);
            },
            $text
        );
        
    }

}