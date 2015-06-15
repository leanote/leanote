var MSG = {"a":"a","aboutMe":"A propos de moi","author":"Auteur","baseInfoSet":"Information de base","blogClass":"Catégorie","blogDesc":"Description","blogLogo":"Logo","blogLogoTips":"Téléverser une image pour remplacer le titre du blog","blogName":"Titre","blogNav":"Navigation du blog","blogNavs":"Navigation","blogSet":"Configuration du blog","cancel":"Annuler","chooseComment":"Système de commentaires","chooseReason":"请选择举报理由","comment":"Commenter","commentSet":"Commentaires","comments":"Commentaires","community":"Communauté","confirm":"Confirmer","confirmDeleteComment":"Êtes-vous sûr?","createdTime":"Créé à","daysAgo":"Il y a plusieurs jours","delete":"Effacer","disqusHelp":"Veuillez renseigner votre identifiant Disqus","domain":"Domaine personnalisé","domainSet":"Domaine","elegant":"Elegant","error":"Erreur","fullBlog":"Blog entier","home":"Accueil","hoursAgo":"Il y a plusieurs heures","justNow":"Seulement maintenant","latestPosts":"Dernières publications","like":"J'aime","minutesAgo":"Il y a quelques minutes","monthsAgo":"Il y a plusieurs mois","more":"Plus...","moreShare":"Plus","navFixed":"Navigation figée sur le côté gauche","needHelp":"Besoin d'aide?","next":"Suivant","noBlog":"Aucun blog","noTag":"Aucune étiquette","none":"Aucun","openComment":"Ouvrir un commentaire?","other":"Autre","previous":"Précédent","qqZone":"QQ Zone","quickLinks":"Liens rapides","renren":"Renren","reply":"Répondre","report":"Rapport","reportBlog?":"举报该博客?","reportComment?":"举报该评论?","reportReason":"Raison","reportReason1":"不友善内容","reportReason2":"广告等垃圾信息","reportReason3":"违法违规内容","reportReason4":"不宜公开讨论的政治内容","reportSuccess":"举报成功, 我们处理后会通知作者, 感谢您的监督","saveSuccess":"Sauvegarde réussie","scanQRCode":"Ouvrez Weichat et scannez le QR Code","signIn":"S'identifier","signUp":"S'incrire","sinaWeibo":"Weibo","subDomain":"Sous-domaine","submitComment":"Soumettre","tencentWeibo":"Tencent Weibo","theme":"Thème","themeSet":"Thème","unlike":"Je n'aime plus","updatedTime":"Mis à jour à","viewers":"Lecteurs","weeksAgo":"Il y a plusieurs semaines","weixin":"Weichat"};
function getMsg(key, data) {
	var msg = MSG[key]
	if(msg) {
		if(data) {
			if(!isArray(data)) {
				data = [data];
			}
			for(var i = 0; i < data.length; ++i) {
				msg = msg.replace("%s", data[i]);
			}
		}
		return msg;
	}
	return key;
}