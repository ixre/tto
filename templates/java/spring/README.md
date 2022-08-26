# 数据查询组件

文件(JAVA版):`ReportComponent.java`

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

文件(Kotlin版):`ReportComponent.kt`

```kotlin
import net.fze.common.Standard
import net.fze.extras.report.*
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.stereotype.Component
import java.sql.Connection
import javax.sql.DataSource

@Component
class ReportComponent : IDbProvider {
    private val ReportHubMap: MutableMap<String, ReportHub> = mutableMapOf()
    private val rootPath = "/query"

    @Inject
    private var ds: DataSource? = null

    override fun getDB(): Connection {
        return try {
            ds!!.connection
        } catch (ex: Exception) {
            ex.printStackTrace()
            throw Error(ex.message)
        }
    }

    private fun lazyInit() {
        ReportHubMap["default"] = ReportHub(
            this,
            "$rootPath/default@query", !Standard.dev()
        )
    }

    private fun getHub(key: String): ReportHub? {
        if (this.ReportHubMap.isEmpty()) {
            this.lazyInit()
        }
        if (key.isEmpty()) return ReportHubMap["default"]
        return ReportHubMap[key]
    }

    fun parseParams(params: String): Params {
        return ReportUtils.parseParams(params)
    }

    fun fetchData(key: String, portal: String, params: Params, page: String, rows: String): DataResult {
        val hub = this.getHub(key) ?: throw Exception("datasource not exists")
        return hub.fetchData(portal, params, page, rows)
    }
}
```