[jQuery Validation Plugin](http://jqueryvalidation.org/) - Form validation made easy
================================

[![Build Status](https://secure.travis-ci.org/jzaefferer/jquery-validation.png)](http://travis-ci.org/jzaefferer/jquery-validation)
[![devDependency Status](https://david-dm.org/jzaefferer/jquery-validation/dev-status.png?theme=shields.io)](https://david-dm.org/jzaefferer/jquery-validation#info=devDependencies)

The jQuery Validation Plugin provides drop-in validation for your existing forms, while making all kinds of customizations to fit your application really easy.

## [Help the project](http://pledgie.com/campaigns/18159)

[![Help the project](http://www.pledgie.com/campaigns/18159.png?skin_name=chrome)](http://pledgie.com/campaigns/18159)

This project is looking for help! [You can donate to the ongoing pledgie campaign](http://pledgie.com/campaigns/18159)
and help spread the word. If you've used the plugin, or plan to use, consider a donation - any amount will help.

You can find the plan for how to spend the money on the [pledgie page](http://pledgie.com/campaigns/18159).

## Getting Started

### Downloading the prebuilt files

Prebuilt files can be downloaded from http://jqueryvalidation.org/

### Downloading the latest changes

The unreleased development files can be obtained by:

 1. [Downloading](https://github.com/jzaefferer/jquery-validation/archive/master.zip) or Forking this repository
 2. [Setup the build](CONTRIBUTING.md#build-setup)
 3. Run `grunt` to create the built files in the "dist" directory

### Including it on your page

Include jQuery and the plugin on a page. Then select a form to validate and call the `validate` method.

```html
<form>
	<input required>
</form>
<script src="jquery.js"></script>
<script src="jquery.validate.js"></script>
<script>
$("form").validate();
</script>
```

Alternatively include jQuery and the plugin via requirejs in your module.

```js
define(["jquery", "jquery.validate"], function( $ ) {
	$("form").validate();
});
```

For more information on how to setup a rules and customizations, [check the documentation](http://jqueryvalidation.org/documentation/).

## Reporting an Issue

**IMPORTANT NOTE ABOUT EMAIL VALIDATION**. As of version 1.12.0 this plugin is using the same regular expression that the [HTML5 specification suggests for browsers to use](http://www.whatwg.org/specs/web-apps/current-work/multipage/states-of-the-type-attribute.html#e-mail-state-%28type=email%29). We will follow their lead and use the same check. If you think the specification is wrong, please report the issue to them.

1. Make sure the problem you're addressing is reproducible.
2. Use http://jsbin.com or http://jsfiddle.net to provide a test page.
3. Indicate what browsers the issue can be reproduced in. **Note: IE Compatibilty modes issues will not be addressed.**
4. What version of the plug-in is the issue reproducible in. Is it reproducible after updating to the latest version.

## Contributing

See the [Contributing Guide](CONTRIBUTING.md)

## License
Copyright (c) 2013 Jörn Zaefferer
Licensed under the MIT license.
