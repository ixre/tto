# 数据查询组件

文件:`ReportDataSource.java`

```java
import net.fze.common.Standard;
import net.fze.extras.report.DataResult;
import net.fze.extras.report.ReportHub;
import net.fze.extras.report.Params;

@Component
public class ReportDataSource {
    private final HashMap<String, ReportHub> ReportHubMap = new HashMap<>();

    @Inject
    private DataSource ds;

    public Connection getDB(String key) {
        try {
            switch(key){
                default: 
                 return this.ds.getConnection();
            }
        } catch (SQLException e) {
            e.printStackTrace();
        }
        throw new Error("can't get any connection");
    }

    private void lazyInit() {
        boolean cache = !Standard.dev();
        ReportHubMap.put("default", new ReportHub( ()->getDB("default"),"/query",cache ));
    }

    private ReportHub getHub(String key) {
        if (this.ReportHubMap.isEmpty())this.lazyInit();
        if(ReportHubMap.containsKey(key))return ReportHubMap.get(key);
        return ReportHubMap.get("default");
    }

    public DataResult fetchData(String key, String portal, Params params, String page, String rows) {
        ReportHub hub = this.getHub(key);
        if (hub == null)throw new Error("datasource not exists");
        return hub.fetchData(portal, params, page, rows);
    }
}
```
