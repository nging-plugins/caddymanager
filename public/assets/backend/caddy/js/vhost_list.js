$(function(){
	App.bindSwitch('input.switch-disabled','click','caddy/vhost_edit');
	var wp = function (raw, xy) {
		return '<div class="wrap">'+App.htmlEncode(raw)+'</div>';
	}, wp2 = function (raw, xy) {
		return '<div class="wrap" style="min-width:100px">'+App.htmlEncode(raw)+'</div>';
	}, genTable = function(list){
		return App.genTable(list, {
			'StatusCode': function (raw, xy) {
                var text=HTTPStatusText(raw);
				return '<span class="label label-' + App.httpStatusColor(raw) + '" title="'+text+'" onclick="$(this).next(\'small\').toggleClass(\'hide\')">' + raw + '</span> <small class="hide">'+text+'</small>';
			},
			'Brower': wp,
			'UserAgent': wp2,
			'URI': wp2,
			'Referer': wp2,
			'': function (raw, xy) {
				return App.htmlEncode(raw);
			}
		})
	};
	$('a.access-log-pretty-show').on('click',function(){
		App.logShow(this,true,'access',genTable);
	})
	App.logShow('#caddy-log-show');
});