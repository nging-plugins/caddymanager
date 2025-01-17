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
  var tmpl=$('#tmplAddVariableRow').html();
  tmpl=tmpl.replace(/\{k\}/g,k);
  tmpl=tmpl.replace(/\{v\}/g,v);
  $(obj).parent().before(tmpl);
}
function addKs(obj,k,v){
  var tmpl=$('#tmplAddVariableRowSingleCell').html();
  tmpl=tmpl.replace(/\{k\}/g,k);
  $(obj).parent().before(tmpl);
}
function initReqPlaceholdersModal(){
  $('#request-placeholders-modal').find('.modal-body').css({"max-height":$(window).height()-200});
}
function copyFormHTML(boxElem){
  var base = $(boxElem).children('.fieldset:first');
  var closeBtn = '<a href="javascript:;" class="label label-danger extra-page-remove" onclick="removeFormHTML(this);" data-toggle="tooltip" title="'+App.t('删除')+'"><i class="fa fa-times"></i></a>';
  var copied = base.clone();
  copied.prepend(closeBtn);
  copied.find('input[type="text"]').val('');
  copied.find('textarea').text('');
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
});