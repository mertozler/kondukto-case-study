# Kondukto Case Study
![Kondukto Logo](https://assets-global.website-files.com/5fec9210c1841a6c20c6ce81/60f4b69942eda3f754b4f2a5_KonduktoLogo.png)


Bu proje bir REST API'sidir. /api/v1/newscan adresine json body olarak gönderilen python github repolarını Bandit ile analiz edip, güvenlik açıklarını analiz eder. Bu analiz datalarını Redis'te depolar ve  scan_id üzerinden bu verilere erişilebilir..

## Kurulum

Projeyi localinizde çalıştırmak için ilk olarak Redis'i docker üzerinde ayağa kaldırmanız gerekmektedir. Unutmayınız: config dosyasında yapmanız gereken değişiklikler olup olmadığını teyit edin ve projeyi çalıştırın.

```bash
docker-compose up -d
```

## Alternatif Yöntemler

Projeyi geliştirirken alternatif yöntemler üstüne düşündüm. Şu anda gelen her istekte Docker üzerinde Bandit'i ayağa kaldırıp bu verileri analiz ediyoruz. Her istekte yeni docker kaldırmamızın sebebi bandit'in container ayağa kaldırırken PWD değişkenini almasından kaynaklı. Buna alternatif bir çözüm bulamadığım için tek bir container'ı arkada çalıştırıp gelen requestler üzerinden scan datalarını almayı düşündüm fakat bunu yapamadım. Bandit'i Docker'da ayağa Ardından analiz çıktılarını kullanıcıya dönüyoruz ve bu haliyle response time'ı artırıyor. Kapsam dışı olsa da biz kullanıcıdan aldığımız her isteği Kafka ile bir diğer servise (bandit-service) scan_id, repo_url payload'ı ile gönderip ardından analiz sonuçlarını kullanıcıya mail olarak dönebiliriz. Böylelikle response time'ı artabilir. İkincil yöntem olarak; kullanıcının isteğini aldıktan sonra scan_id, repo_url şeklinde Redis DB'ye kaydeder ve response döneriz. Bir background job belirli aralıklarla henüz analiz edilmemiş repoları arkada analiz eder ve sonuçları kullanıcıya mail olarak döner. Bu projede kapsam dışı olduğundan bunları geliştirmedim. 

## Proje Nasıl İlerletildi?
Proje github üzerinden ilerletildi. Geliştirmeyi daha yakından takip edebilmek ve hızlı geliştirmek için Agile metodolojisini kendimce uygulamak için Backloglar oluşturdum. Buradan erişim sağlayabilirsiniz; [tıklayın](https://docs.google.com/spreadsheets/d/1-NHtcoOE55HmqCgnzKbYRcq0llRjUcl2_nc3n8ApI5c/edit?usp=sharing)
