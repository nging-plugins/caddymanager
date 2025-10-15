var inputRequestPlaceholder=null;
function showRequestPlaceholders(a){
  if(a==null) a='#rewrite_rule';
  inputRequestPlaceholder=a;
  $('#request-placeholders-modal').niftyModal('show');
}
function insertAdvRewrite(){
  App.insertAtCursor($('#rewrite_rule')[0],'rewrite / {\n\
    regexp pattern\n\
    ext    extensions...\n\
    if     a cond b\n\
    if_op  [and|or]\n\
    to     destinations...\n\
}\n',23,31);
}
function insertAdvRedir(){
  App.insertAtCursor($('#rewrite_rule')[0],'redir 301 {\n\
    if {>X-Forwarded-Proto} is http\n\
    /  https://{host}{uri}\n\
}\n');
}
function insertCondForRewrite(obj){
  App.insertAtCursor($('#rewrite_rule')[0],$(obj).children('strong').text()+' ');
}
function addKVs(obj,k,v){
  var tmpl=$('#tmplAddVariableRow').html(),namePrefix=$(obj).closest('table').data('nameprefix');
  if(namePrefix){
    var newName=formInputNameWithPrefix(k,namePrefix);
    if(newName!=='') k=newName;
  }
  tmpl=tmpl.replace(/\{k\}/g,k);
  tmpl=tmpl.replace(/\{v\}/g,v);
  $(obj).parent().before(tmpl);
}
function addKs(obj,k){
  var tmpl=$('#tmplAddVariableRowSingleCell').html(),namePrefix=$(obj).closest('table').data('nameprefix');
  if(namePrefix){
    var newName=formInputNameWithPrefix(k,namePrefix);
    if(newName!=='') k=newName;
  }
  tmpl=tmpl.replace(/\{k\}/g,k);
  $(obj).parent().before(tmpl);
}
function initReqPlaceholdersModal(){
  $('#request-placeholders-modal').find('.modal-body').css({"max-height":$(window).height()-200});
}
function formInputNameWithPrefix(name,namePrefix){
  if(name.startsWith(namePrefix)) return '';
  if(name.endsWith(']')){
    var pos=name.indexOf('['),newName='';
    if(pos>0){
      newName=namePrefix+'['+name.substring(0,pos)+']'+name.substring(pos);
    }else if(pos==0){
      newName=namePrefix+name;
    }else{
      newName=namePrefix+'['+name;
    }
  }else{
    newName=namePrefix+'['+name+']';
  }
  return newName;
}
function addNamePrefix(a,namePrefix){
  var name=$(a).attr('name'), newName=formInputNameWithPrefix(name,namePrefix);
  if(newName==='') return;
  $(a).attr('name',newName);
  var id=$(a).attr('id');
  if(!id) return;
  $(a).attr('id',namePrefix+id);
  $(a).next('label[for="'+id+'"]').attr('for',namePrefix+id);
}
function copyFormHTML(boxElem,namePrefix,noClean){
  var base = $(boxElem).children('.fieldset:first');
  var closeBtn = '<a href="javascript:;" class="label label-danger extra-page-remove" onclick="removeFormHTML(this);" data-toggle="tooltip" title="'+App.t('删除')+'"><i class="fa fa-times"></i></a>';
  var copied = base.clone();
  copied.prepend(closeBtn);
  if(!noClean){
    copied.find('input[type="text"]').val('');
    copied.find('textarea').text('');
  }
  if(copied[0].hasAttribute('id')) copied[0].removeAttribute('id');
  if(namePrefix) {
    var indexName=namePrefix+'_index[]',index=$(boxElem).find('input[name="'+indexName+'"]:last').val()||0;
    index++;
    namePrefix=namePrefix+'['+index+']';
    copied.find('[name]').each(function(){addNamePrefix(this,namePrefix)});
    copied.find('.form-group table').attr('data-nameprefix',namePrefix);
    copied.prepend('<input type="hidden" name="'+indexName+'" value="'+index+'" />');
  }
  if(base.next('.fieldset').length>0){
    base.siblings('.fieldset:last').after(copied);
  }else{
    base.after(copied);
  }
}
function removeFormHTML(obj){
  $(obj).parent('.fieldset').remove();
}
function showEngineElements($form){
  var engine=$('#serverIdent').children('option:selected').data('engine');
  $form.find('.engine-show-'+engine).show();
  $form.find('.engine-hide-'+engine).hide();
  $form.find('[class^="engine-hide"]:not(.engine-hide-'+engine+')').show();
  $form.find('[class^="engine-show"]:not(.engine-show-'+engine+')').hide();
}
$(function(){
initReqPlaceholdersModal();
$(window).on('resize',function(){
  initReqPlaceholdersModal();
});
$('#vhostForm [data-slide-settings]').each(function(){
  var eid='#'+$(this).attr('name')+'-settings';
  if($(eid).data('history'))$(eid).data('history',null);
});
function changeTabsTitleStatus(tabPanel,show){
  if(tabPanel.length<1) return;
  var tabPanelId = tabPanel.attr('id');
  if(!tabPanelId) return;
  var tabLink=$('.nav-tabs>li>a[href="#'+tabPanelId+'"]');
  if(tabLink.length<1) return;
  var tabItem=$('.nav-tabs>li>a[href="#'+tabPanelId+'"]').parent('li');
  if(show){
    tabItem.removeClass('not-enabled');
  }else if(!tabItem.hasClass('not-enabled')){
    tabItem.addClass('not-enabled');
  }
}
var requiredWWWRootPathFields=['fastcgi','browse'];//必须设置“网站位置”的模块
$('#vhostForm [data-slide-settings]').on('click',function(){
  var slide=$(this).data('slide-settings'),name=$(this).attr('name');
  if(slide=="")slide=$(this).val()=='1'?'show':'hide';
  var eid='#'+name+'-settings',innerid='#'+name+'-settings-inner';
  if($(innerid).length<1) innerid=eid;
  var history=$.trim($(innerid).html());
  for(var i=0;i<requiredWWWRootPathFields.length;i++){
    if(name==requiredWWWRootPathFields[i]){
      if(slide=='show'){
        $('#www-root-path').prop('required',true);
      }else{
        var checked=false;
        for(var j=0;j<requiredWWWRootPathFields.length;j++){
          var f=requiredWWWRootPathFields[j];
          if($('#'+f+'-enable').prop('checked')){
            checked=true;
            break;
          }
        }
        $('#www-root-path').prop('required',checked);
      }
      break;
    }
  }
  if(slide=='show'){
    if(!history){
      if($(innerid).data('history')){
        $(innerid).html($(innerid).data('history'));
        App.passwordInputShowPassword(innerid);
        showEngineElements($(innerid));
      }else{
        $.get(BACKEND_URL+'/caddy/addon_form',{addon:name},function(r){
          $(innerid).html(r);
          App.passwordInputShowPassword(innerid);
          showEngineElements($(innerid));
        },'html');
      }
    }else{
      showEngineElements($(innerid));
    }
    $(eid).removeClass('hide');
  }else{
    if(history)$(innerid).data('history',history);
    $(innerid).empty();
    $(eid).addClass('hide');
  }
  var tabPanel=$(this).closest('.tab-pane');
  changeTabsTitleStatus(tabPanel,$(this).val()=='1');
});
$('#vhostForm [data-slide-settings="show"]:checked').trigger('click');
$('#request-placeholders-modal .modal-body table tbody td b').off().on('click',function(){
  App.insertAtCursor($(inputRequestPlaceholder)[0],$(this).text());
  $('#request-placeholders-modal').niftyModal('hide');
});
$('[data-popover="popover"]').on('click',function(){
  $(this).siblings('[data-popover="popover"]').popover('hide');
});
$('#serverIdent').on('change',function(){
  var $form=$('#vhostForm');
  showEngineElements($form);
}).trigger('change');
App.searchFS('#www-root-path',20,'dir');
$('#vhostForm').on('change','#tls_acme_dns_provider',function(){
  var provider=$(this).val(),$inputs=$('#tls-dns-provider-inputs');
  if(!provider){
    $inputs.empty();
    return;
  }
  var dataInputs=$(this).children('option:selected').data('inputs');
  var params=[];
  for(var i=0;i<dataInputs.length;i++){
    var input=dataInputs[i];
    params.push({Label:App.t(input.Label),Input:input.Input,Help:input.Help})
  }
  if($inputs.length<1){
    $inputs=$('<div id="tls-dns-provider-inputs"></div>');
    $(this).closest('.form-group').after($inputs);
  }
  $inputs.html(template('tpl-tls-dns-provider-input',{inputs:params}));
}).trigger('change');
$('#vhostForm').on('change','#tls_ca_name',function(){
  var ca=$(this).val(),$inputs=$('#tls-ca-inputs');
  if(!ca){
    $inputs.empty();
    return;
  }
  var $sel=$(this).children('option:selected'),needEAB=$sel.data('need-eab'),needToken=$sel.data('need-token');
  if($inputs.length<1){
    $inputs=$('<div id="tls-ca-inputs"></div>');
    $(this).closest('.form-group').after($inputs);
  }else{
    $inputs.empty();
  }
  if(needEAB) $inputs.append(template('tpl-tls-ca-eab',{}));
  if(needToken) $inputs.append(template('tpl-tls-ca-token',{}));
}).trigger('change');
});