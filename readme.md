# 리뷰 분류 프로젝트
- 네이버 블로그에서 상품 리뷰를 보다보면 조회수를 늘리기위해 매크로 형식으로 작성된 글을 볼수 있다. 이러한 글들을 정상적이지 않은 리뷰(이하 스팸리뷰)와 그렇지 않은 리뷰를 구분하는 프로젝트
- 리뷰는 크게 2가지 카테고리로 분류한다. 리뷰, 스팸리뷰 -> 추후 된다면 커미션을 받은 상업적 리뷰도 추가해본다.

- 리뷰의 예시
- 스팸리뷰의 예시

## 프로젝트의 구조
- 데이터를 크롤링 하는 Go프로그램
- 크롤링한 데이터를 적재하는 Elastic Search
- 각종 설정과 그 외 잡다한 설정을 적재하는 MariaDB
- 분류(딥러닝을 통한)를 하기위한 python
- 분류의 결과를 표시해줄 웹(Srping boot)

### 현재 진행상황
1. Go를 사용해서 네이버 블로그의 데이터 가져오기 (현재 진행)
2. Docker를 이용해 Elastic Search와 kibana 구축 

### Go를 이용한 스크래퍼 설정
- 크롤은 현재 네이버 블로그만 생각하고 있지만 추후 수집하는 매체가 늘어날 수 있으므로 여러 PC에서 수집할 필요성이 존재(수집을 빠르게 돌리면 사이트에서 ip를 차단하므로 느림)
- 이러한 문제를 해결하기위해 각 사이트별로 별도의 프로그램을 작성하여 cron을 사용하여 수집. 수집의 로그는 MariaDB에 적재하는 방식으로 수집 내역을 관리
- HTML Document를 읽기 파싱하기 위한 GoQuery 라이브러리설치 -> go get github.com/PuerkitoBio/goquery
- naver 블로그 내역을 제공하는 API가 존재하여 API로 블로그 내역을 수집
    - 제공되는 데이터는 XML과 JSON 타입으로 제공되며 제공되는 데이터의 내역은
        -title, link, discription, blogername, blogerlink, postdate
    - title, link, discription, postdate를 사용할 예정
- 네이버 개발자센터 링크: <https://developers.naver.com/docs/serviceapi/search/blog/blog.md#%EB%B8%94%EB%A1%9C%EA%B7%B8>